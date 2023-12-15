package challenge

import (
	"strconv"
	"strings"
)

func Challenge(input string) string {
	items := strings.Split(input, ",")

	tally := 0

	for _, item := range items {
		tally += hashAlgorithm(item)
	}

	return strconv.Itoa(tally)
}

func hashAlgorithm(input string) int {
	var val int

	for _, c := range input {
		if c == '\n' {
			continue
		}
		val += int(c)
		val *= 17
		val = val % 256
	}

	return val
}
