package hyphenate

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestSplitWord(t *testing.T) {
	d, err := LoadDictionary("hyph_en_US.dic")
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Open("test-cases.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	count := 0

	for scanner.Scan() {
		pieces := strings.Split(scanner.Text(), " ")
		original := pieces[0]
		expected := pieces[1]
		actual := d.Hyphenate(original, "-")

		if actual != expected {
			t.Errorf("Expected %q -> %q, not %q", original, expected, actual)
			count++
		}
	}

	if count > 0 {
		t.Errorf("%d words were incorrect", count)
	}
}

func BenchmarkLoadDictionary(b *testing.B) {
	stat, err := os.Stat("hyph_en_US.dic")
	if err != nil {
		b.Fatal(err)
	}

	b.SetBytes(stat.Size())
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		LoadDictionary("hyph_en_US.dic")
	}
}

func BenchmarkLoadDictionaryGzip(b *testing.B) {
	stat, err := os.Stat("hyph_en_US.dic.gz")
	if err != nil {
		b.Fatal(err)
	}

	b.SetBytes(stat.Size())
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		LoadDictionary("hyph_en_US.dic.gz")
	}
}

func BenchmarkHyphenate(b *testing.B) {
	d, err := LoadDictionary("hyph_en_US.dic")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Hyphenate("electroencephalograph's", "-")
	}
}
