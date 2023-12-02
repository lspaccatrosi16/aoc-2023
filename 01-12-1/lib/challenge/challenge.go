package challenge

import (
	"fmt"
	"strconv"
	"strings"
)

func Challenge(input string) string {
	lines := strings.Split(input, "\n")
	tally := 0

	for ln, line := range lines {
		var first, last rune
		for i := 0; i < len(line); i++ {
			if line[i] >= '0' && line[i] <= '9' {
				first = rune(line[i])
				break
			}
		}
		for i := len(line) - 1; i >= 0; i-- {
			if line[i] >= '0' && line[i] <= '9' {
				last = rune(line[i])
				break
			}
		}

		if first == 0 || last == 0 {
			msg := fmt.Errorf("line %d: does not have at least one number", ln+1)
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
