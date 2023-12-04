package challenge

import (
	"fmt"
	"strconv"
	"strings"
)

var exists = struct{}{}

func Challenge(input string) string {
	lines := strings.Split(input, "\n")

	tally := 0

	for i, l := range lines {
		tally += handleLine(l, i)
	}

	return strconv.Itoa(tally)
}

func handleLine(line string, num int) int {
	components := strings.Split(line, ":")

	if len(components) != 2 {
		msg := fmt.Errorf("line %d has an incorrect format (components)", num+1)
		panic(msg)
	}

	lists := strings.Split(components[1], "|")

	if len(lists) != 2 {
		msg := fmt.Errorf("line %d has an incorrect format (lists)", num+1)
		panic(msg)
	}

	winning := strings.Split(strings.Trim(lists[0], " "), " ")
	given := strings.Split(strings.Trim(lists[1], " "), " ")

	winningMap := map[int]interface{}{}

	for _, n := range winning {
		if n == "" {
			continue
		}

		vInt, err := strconv.Atoi(n)

		if err != nil {
			msg := fmt.Errorf("line %d has an incorrect format (parseint): '%s'", num+1, n)
			panic(msg)
		}

		winningMap[vInt] = exists
	}

	numMatches := 0

	for _, n := range given {
		if n == "" {
			continue
		}

		vInt, err := strconv.Atoi(n)

		if err != nil {
			msg := fmt.Errorf("line %d has an incorrect format (parseint): '%s'", num+1, n)
			panic(msg)
		}

		_, exists := winningMap[vInt]
		if exists {
			numMatches++
		}
	}

	var score int

	if numMatches == 0 {
		score = 0
	} else if numMatches == 1 {
		score = 1
	} else {
		score = 1
		for i := 1; i < numMatches; i++ {
			score *= 2
		}
	}

	fmt.Printf("Card %d: %d Matches => Score %d\n", num+1, numMatches, score)
	return score
}
