package challenge

import (
	"fmt"
	"strconv"
	"strings"
)

type ScratchCard struct {
	CardNo   int
	Quantity int
	Matches  int
	Contents string
}

func (s *ScratchCard) Eval() int {
	if s.Matches >= 0 {
		return s.Matches
	} else {
		matches := handleLine(s.Contents, s.CardNo)
		s.Matches = matches
		return matches
	}
}

type ScractchCardList map[int]*ScratchCard

func (s ScractchCardList) Count() int {
	tally := 0
	for _, v := range s {
		tally += v.Quantity
	}
	return tally
}

var exists = struct{}{}

func Challenge(input string) string {
	lines := strings.Split(input, "\n")

	list := ScractchCardList{}

	for i, l := range lines {
		list[i] = &ScratchCard{
			CardNo:   i,
			Quantity: 1,
			Matches:  -1,
			Contents: l,
		}
	}

	queue := []*ScratchCard{}

	for _, v := range list {
		queue = append(queue, v)
	}

	for i := 0; ; i++ {
		if i >= len(queue) {
			break
		}

		matches := queue[i].Eval()
		for j := 1; j <= matches; j++ {
			queue[i].Quantity++
			queue = append(queue, list[queue[i].CardNo+j])
		}
	}

	return strconv.Itoa(list.Count())
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

	return numMatches
}
