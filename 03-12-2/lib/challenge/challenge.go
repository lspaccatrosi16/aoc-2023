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
	gearMap := map[int]map[int][]int{}

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

					isPartNum, x, y := checkAdjacent(start, end, curLine, aboveLine, belowLine)
					y = i + y

					if isPartNum {
						var arr []int

						if vx, existsx := gearMap[x]; existsx {
							if vy, existsy := vx[y]; existsy {
								arr = vy
							}
						} else {
							gearMap[x] = map[int][]int{}
						}

						arr = append(arr, parsedNum)
						gearMap[x][y] = arr
					}

					curNum = ""
				}
			}
		}

	}
	for _, vx := range gearMap {
		for _, vy := range vx {
			if len(vy) != 2 {
				continue
			}

			ratio := vy[0] * vy[1]
			tally += ratio

		}
	}

	return strconv.Itoa(tally)
}

// bool, x coord, y coord
func checkAdjacent(start, end int, currentL, aboveL, belowL string) (bool, int, int) {

	// Check left
	if start > 0 {
		if isGear(currentL[start-1]) {
			return true, start - 1, 0
		}

		if aboveL != "" {
			if isGear(aboveL[start-1]) {
				return true, start - 1, -1
			}
		}

		if belowL != "" {
			if isGear(belowL[start-1]) {
				return true, start - 1, 1
			}
		}
	}

	// Check above / below
	for i := start; i <= end; i++ {
		if aboveL != "" {
			if isGear(aboveL[i]) {
				return true, i, -1
			}
		}

		if belowL != "" {
			if isGear(belowL[i]) {
				return true, i, 1
			}
		}
	}

	// Check right
	if end+1 < len(currentL) {
		if isGear(currentL[end+1]) {
			return true, end + 1, 0
		}

		if aboveL != "" {
			if isGear(aboveL[end+1]) {
				return true, end + 1, -1
			}
		}

		if belowL != "" {
			if isGear(belowL[end+1]) {
				return true, end + 1, 1
			}
		}
	}
	return false, 0, 0
}

func isNum(b byte) bool {
	return b >= '0' && b <= '9'
}

func isGear(b byte) bool {
	return b == '*'
}
