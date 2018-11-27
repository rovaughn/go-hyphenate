import pyphen

d = pyphen.Pyphen("hyph_en_US.dic")

with open("/usr/share/dict/american-english") as words:
    with open("test-cases.txt", "w") as out:
        for word in words:
            word = word.strip()
            assert " " not in word and "-" not in word
            print(word, d.inserted(word), file=out)
