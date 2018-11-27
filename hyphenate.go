package hyphenate

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type pattern struct {
	offset int
	values []int
}

type Dictionary struct {
	patterns map[string]pattern
}

func (d *Dictionary) addPattern(p string) {
	value := 0
	values := make([]int, 0)
	tags := make([]rune, 0)
	for _, r := range p {
		if '0' <= r && r <= '9' {
			value = int(r - '0')
		} else {
			values = append(values, value)
			tags = append(tags, r)
			value = 0
		}
	}

	values = append(values, value)

	allZero := true
	for _, value := range values {
		if value != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		return
	}

	offset := 0
	for values[offset] == 0 {
		offset++
	}

	end := len(values)
	for values[end-1] == 0 {
		end--
	}

	d.patterns[string(tags)] = pattern{
		offset: offset,
		values: values[offset:end],
	}
}

func LoadDictionary(filename string) (*Dictionary, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d := &Dictionary{
		patterns: make(map[string]pattern),
	}

	scanner := bufio.NewScanner(f)

	if !scanner.Scan() {
		return nil, fmt.Errorf("Expected encoding as first line")
	}

	encoding := scanner.Text()
	if encoding != "UTF-8" {
		return nil, fmt.Errorf("Unknown encoding %q", encoding)
	}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "LEFTHYPHENMIN") || strings.HasPrefix(line, "RIGHTHYPHENMIN") || strings.HasPrefix(line, "COMPOUNDLEFTHYPHENMIN") || strings.HasPrefix(line, "COMPOUNDRIGHTHYPHENMIN") {
			continue
		}
		d.addPattern(line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return d, nil
}

// NOTE we can cache the length of the longest key in patterns; no need to look
// for a pattern longer than that.
func (d *Dictionary) getPositions(word string) []int {
	word = strings.ToLower(word)
	pointedWord := []rune("." + word + ".")
	references := make([]int, len(pointedWord)+1)

	for i := 0; i < len(pointedWord); i++ {
		for j := i + 1; j <= len(pointedWord); j++ {
			slice := pointedWord[i:j]
			pattern, ok := d.patterns[string(slice)]
			if !ok {
				continue
			}

			for k, value := range pattern.values {
				if value > references[i+pattern.offset+k] {
					references[i+pattern.offset+k] = value
				}
			}
		}
	}

	points := make([]int, 0, len(references)/2)

	for i, reference := range references {
		point := i - 1
		if point >= 2 && point <= len(pointedWord)-4 && reference%2 == 1 {
			points = append(points, point)
		}
	}

	return points
}

func (d *Dictionary) SplitWord(word string) []string {
	points := d.getPositions(word)

	if len(points) == 0 {
		return []string{word}
	}

	runes := []rune(word)

	result := make([]string, 0)

	result = append(result, string(runes[:points[0]]))
	for i := 0; i < len(points)-1; i++ {
		result = append(result, string(runes[points[i]:points[i+1]]))
	}
	result = append(result, string(runes[points[len(points)-1]:]))
	return result
}

func (d *Dictionary) Hyphenate(word string, hyphen string) string {
	return strings.Join(d.SplitWord(word), hyphen)
}
