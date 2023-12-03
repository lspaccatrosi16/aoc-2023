package challenge

import (
	"fmt"
	"strconv"
	"strings"
)

func Challenge(input string) string {
	lines := strings.Split(input, "\n")

	tally := 0

	var curLine, aboveLine, belowLine, curNum string

	for i := 0; i < len(lines); i++ {
		if i > 0 {
			aboveLine = lines[i-1]
		}
		curLine = lines[i]
		if i+1 < len(lines) {
			belowLine = lines[i+1]
		}

		for j := 0; j < len(curLine); j++ {
			if isNum(curLine[j]) {
				curNum += string(curLine[j])
				if j+1 < len(curLine) && isNum(curLine[j+1]) {
					continue
				} else {
					parsedNum, err := strconv.Atoi(curNum)
					if err != nil {
						msg := fmt.Errorf("line %d: error parsing part number at position [%d:%d]", i+1, j+1-len(curNum), j+1)
						panic(msg)
					}

					start := j - len(curNum) + 1
					end := j

					if start < 0 || end < 0 {
						msg := fmt.Errorf("line %d: start, end out of range [%d:%d] %d, %d. curNum: '%s'", i+1, j+1-len(curNum), j+1, start, end, curNum)
						panic(msg)
					}

					isPartNum := checkAdjacent(start, end, curLine, aboveLine, belowLine)

					if isPartNum {
						tally += parsedNum
					}

					curNum = ""
				}
			}
		}
	}

	return strconv.Itoa(tally)
}

func checkAdjacent(start, end int, currentL, aboveL, belowL string) bool {

	// Check left
	if start > 0 {
		if isSym(currentL[start-1]) {
			return true
		}

		if aboveL != "" {
			if isSym(aboveL[start-1]) {
				return true
			}
		}

		if belowL != "" {
			if isSym(belowL[start-1]) {
				return true
			}
		}
	}

	// Check above / below
	for i := start; i <= end; i++ {
		if aboveL != "" {
			if isSym(aboveL[i]) {
				return true
			}
		}

		if belowL != "" {
			if isSym(belowL[i]) {
				return true
			}
		}
	}

	// Check right
	if end+1 < len(currentL) {
		if isSym(currentL[end+1]) {
			return true
		}

		if aboveL != "" {
			if isSym(aboveL[end+1]) {
				return true
			}
		}

		if belowL != "" {
			if isSym(belowL[end+1]) {
				return true
			}
		}
	}
	return false
}

func isNum(b byte) bool {
	return b >= '0' && b <= '9'
}

func isSym(b byte) bool {
	return (b >= '!' && b <= '/' && b != '.') || (b >= ':' && b <= '@') || (b >= '[' && b <= '`') || (b >= '{' && b <= '~')
}
