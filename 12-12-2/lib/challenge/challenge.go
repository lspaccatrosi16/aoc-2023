package challenge

import (
	"strconv"
	"strings"
	"sync"
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

type RecordSet []SpringType

func (r RecordSet) IsPossible(nums ...int) bool {
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

func (r *RecordSet) String() string {
	str := ""
	for _, rec := range *r {
		str += rec.String()
	}
	return str
}

func (r *RecordSet) SplitAtBroken() []RecordSet {
	sets := []RecordSet{}

	var curSet RecordSet

	for i := 0; i < len(*r); i++ {
		curType := (*r)[i]
		switch curType {
		case Broken:
			if len(curSet) > 0 {
				sets = append(sets, curSet)
				curSet = RecordSet{}
			}
		case Operational, Unknown:
			curSet = append(curSet, curType)
		}
	}

	if len(curSet) > 0 {
		sets = append(sets, curSet)
	}

	return sets
}

func Challenge(input string) string {
	lines := strings.Split(input, "\n")

	tally := 0

	wg := &sync.WaitGroup{}

	results := make(chan int)

	for _, l := range lines {
		if l == "" {
			continue
		}
		wg.Add(1)
		go handleLine(wg, l, results)
	}

	go monitor(wg, results)

	for i := range results {
		tally += i
	}

	return strconv.Itoa(tally)
}

func monitor(wg *sync.WaitGroup, results chan int) {
	wg.Wait()
	close(results)
}

func handleLine(wg *sync.WaitGroup, l string, results chan int) {
	defer wg.Done()
	records := RecordSet{}

	components := strings.Split(l, " ")
	unknownIdx := []int{}

	newSprings := []string{}
	newGroups := []string{}

	for i := 0; i < 5; i++ {
		newSprings = append(newSprings, components[0])
		newGroups = append(newGroups, components[1])
	}

	components[0] = strings.Join(newSprings, "?")
	components[1] = strings.Join(newGroups, ",")

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

	output := trySet(records, unknownIdx, groups)
	results <- output
}

func trySet(records RecordSet, unknownIdx []int, groups []int) int {
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
	count += trySet(records, unknownIdx[1:], groups)

	// try broken
	records[thisWork] = Broken
	count += trySet(records, unknownIdx[1:], groups)

	return count
}
