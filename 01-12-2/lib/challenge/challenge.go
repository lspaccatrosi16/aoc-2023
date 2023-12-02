package challenge

import (
	"fmt"
	"strconv"
	"strings"
)

func Challenge(input string) string {
	lines := strings.Split(input, "\n")
	tally := 0

	numberWords := []string{
		"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
	}

	for ln, line := range lines {
		var first, last rune

		curWord := ""

		curRune := rune(0)

		for i := 0; i < len(line); i++ {
			if line[i] >= '0' && line[i] <= '9' {
				curRune = rune(line[i])
			} else {
				found := false

				prov := curWord + string(line[i])

				for _, w := range numberWords {
					if strings.HasPrefix(w, prov) || w == prov {
						curWord = prov
						found = true
						break
					} else if strings.HasPrefix(w, string(line[i])) {
						curWord = string(line[i])
						found = true
					}
				}

				if !found {
					curWord = ""
					continue
				}

				switch curWord {
				case "zero":
					curRune = '0'
				case "one":
					curRune = '1'
				case "two":
					curRune = '2'
				case "three":
					curRune = '3'
				case "four":
					curRune = '4'
				case "five":
					curRune = '5'
				case "six":
					curRune = '6'
				case "seven":
					curRune = '7'
				case "eight":
					curRune = '8'
				case "nine":
					curRune = '9'
				}
			}
			if curRune != 0 {
				first = curRune
				break
			}
		}

		curRune = 0
		curWord = ""

		for i := len(line) - 1; i >= 0; i-- {
			if line[i] >= '0' && line[i] <= '9' {
				curRune = rune(line[i])
			} else {
				found := false

				prov := string(line[i]) + curWord

				for _, w := range numberWords {
					if strings.HasSuffix(w, prov) || w == prov {
						curWord = prov
						found = true
						break
					} else if strings.HasSuffix(w, string(line[i])) {
						curWord = string(line[i])
						found = true
					}
				}

				if !found {
					curWord = ""
					continue
				}

				switch curWord {
				case "zero":
					curRune = '0'
				case "one":
					curRune = '1'
				case "two":
					curRune = '2'
				case "three":
					curRune = '3'
				case "four":
					curRune = '4'
				case "five":
					curRune = '5'
				case "six":
					curRune = '6'
				case "seven":
					curRune = '7'
				case "eight":
					curRune = '8'
				case "nine":
					curRune = '9'
				}
			}
			if curRune != 0 {
				last = curRune
				break
			}
		}

		if first == 0 || last == 0 {
			msg := fmt.Errorf("line %d: does not have at least one number (f: %s, l: %s)", ln+1, string(first), string(last))
			panic(msg)
		}

		num, err := strconv.Atoi(string([]rune{first, last}))

		if err != nil {
			panic(err)
		}

		tally += num
	}

	return strconv.Itoa(tally)
}
