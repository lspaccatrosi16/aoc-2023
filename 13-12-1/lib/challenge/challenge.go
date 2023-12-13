package challenge

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

func Challenge(input string) string {
	results := make(chan int)
	patterns := strings.Split(input, "\n\n")

	wg := &sync.WaitGroup{}

	for _, p := range patterns {
		wg.Add(1)
		go handlePattern(wg, p, results)
	}

	go monitor(wg, results)

	tally := 0

	for res := range results {
		tally += res
	}

	return strconv.Itoa(tally)
}

func monitor(wg *sync.WaitGroup, results chan int) {
	wg.Wait()
	close(results)
}

type PatternAxis []string

func (p PatternAxis) IsSymmetric(startOffset int) bool {
	if startOffset == 0 {
		return false
	}
	if startOffset >= len(p) {
		return false
	}

	// fmt.Println("offset", startOffset)

	var searchStart, searchEnd, itrLen int

	if startOffset*2 > len(p) {
		searchStart = (2 * startOffset) - len(p)
		searchEnd = len(p) - 1
		itrLen = len(p) - startOffset
	} else {
		searchStart = 0
		searchEnd = 2*startOffset - 1
		itrLen = startOffset
	}

	// fmt.Println(searchStart, searchEnd, itrLen)

	for i := 0; i < itrLen; i++ {
		// fmt.Println(searchStart+i, searchEnd-i)
		// fmt.Println(p[searchStart+i], p[searchEnd-i], p[searchStart+i] == p[searchEnd-i])
		if p[searchStart+i] != p[searchEnd-i] {
			return false
		}
	}
	return true
}

func handlePattern(wg *sync.WaitGroup, pattern string, results chan int) {
	defer wg.Done()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR CAUGHT")
			fmt.Println(pattern)
			fmt.Println(r)
			panic("error")
		}
	}()

	if pattern == "" {
		return
	}

	lines := strings.Split(pattern, "\n")
	patRows := make(PatternAxis, len(lines))
	patCols := make(PatternAxis, len(lines[0]))

	for y, l := range lines {
		patRows[y] = l
		for x, c := range l {
			patCols[x] += string(c)
		}
	}

	thisRes := 0

	// fmt.Println("Inspect X")
	for x := 0; x <= len(patCols); x++ {
		if patCols.IsSymmetric(x) {
			thisRes += x
			break
		}
	}

	// fmt.Println("Inspect Y")
	for y := 0; y <= len(patRows); y++ {
		if patRows.IsSymmetric(y) {
			thisRes += y * 100
			break
		}
	}

	if thisRes == 0 {
		fmt.Println("PATTERN MISS: ", "\n", pattern)
	}

	results <- thisRes
}
