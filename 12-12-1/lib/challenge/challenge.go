package challenge

import (
	"strconv"
	"strings"
)

type SpringType int

const (
	Operational SpringType = iota
	Broken
	Unknown
)

func (s SpringType) String() string {
	switch s {
	case Operational:
		return "#"
	case Broken:
		return "."
	case Unknown:
		return "?"
	default:
		return "0"
	}
}

type RecordLine []SpringType

func (r RecordLine) IsPossible(nums ...int) bool {
	foundNums := []int{}
	curCount := 0
	for _, rec := range r {
		if rec == Operational {
			curCount++
		} else if rec == Broken {
			if curCount > 0 {
				foundNums = append(foundNums, curCount)
				curCount = 0
			}
		} else {
			panic("should only recieve operational and broken springs")
		}
	}

	if curCount > 0 {
		foundNums = append(foundNums, curCount)
	}

	if len(nums) != len(foundNums) {
		return false
	}

	for i := 0; i < len(foundNums); i++ {
		if foundNums[i] != nums[i] {
			return false
		}
	}
	return true
}

func Challenge(input string) string {
	lines := strings.Split(input, "\n")

	tally := 0

	for _, l := range lines {
		if l == "" {
			continue
		}
		val := handleLine(l)
		tally += val

	}

	return strconv.Itoa(tally)
}

func handleLine(l string) int {
	records := RecordLine{}

	components := strings.Split(l, " ")
	unknownIdx := []int{}

	for i, c := range components[0] {
		switch c {
		case '#':
			records = append(records, Operational)
		case '.':
			records = append(records, Broken)
		case '?':
			records = append(records, Unknown)
			unknownIdx = append(unknownIdx, i)
		}
	}

	groups := []int{}

	for _, n := range strings.Split(components[1], ",") {
		vInt, err := strconv.Atoi(n)
		if err != nil {
			panic("int conv error")
		}
		groups = append(groups, vInt)
	}

	return tryLine(records, unknownIdx, groups)
}

func tryLine(records RecordLine, unknownIdx []int, groups []int) int {
	if len(unknownIdx) == 0 {
		if records.IsPossible(groups...) {
			return 1
		} else {
			return 0
		}
	}

	thisWork := unknownIdx[0]

	count := 0

	// try operational
	records[thisWork] = Operational
	count += tryLine(records, unknownIdx[1:], groups)

	// try broken
	records[thisWork] = Broken
	count += tryLine(records, unknownIdx[1:], groups)

	return count
}
