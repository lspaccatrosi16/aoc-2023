package challenge

import (
	"fmt"
	"strconv"
	"strings"
)

const MAX_RED = 12
const MAX_GREEN = 13
const MAX_BLUE = 14

func Challenge(input string) string {
	tally := 0

	lines := strings.Split(input, "\n")

	for i, line := range lines {
		tally += handleGame(line, i)
	}

	return strconv.Itoa(tally)
}

func handleGame(line string, n int) int {
	components := strings.Split(line, ":")
	if len(components) != 2 {
		msg := fmt.Errorf("line %d has an invalid format (:)", n+1)
		panic(msg)
	}

	idStr := strings.Split(components[0], " ")

	if len(idStr) != 2 {
		msg := fmt.Errorf("line %d has an invalid format (gameID)", n+1)
		panic(msg)
	}

	id, err := strconv.Atoi(idStr[1])

	if err != nil {
		msg := fmt.Errorf("line %d has an invalid format (parseID)", n+1)
		panic(msg)
	}

	frames := strings.Split(components[1], ";")
	for fN, frame := range frames {
		values := strings.Split(frame, ",")
		for vN, val := range values {
			fields := strings.Split(strings.Trim(val, " "), " ")
			numCube, err := strconv.Atoi(fields[0])

			if err != nil {
				msg := fmt.Errorf("line %d has int parse error (frame %d, val %d): '%s'", n+1, fN+1, vN+1, fields[0])
				panic(msg)
			}

			switch fields[1] {
			case "red":
				if numCube > MAX_RED {
					return 0
				}
			case "green":
				if numCube > MAX_GREEN {
					return 0
				}
			case "blue":
				if numCube > MAX_BLUE {
					return 0
				}
			default:
				msg := fmt.Errorf("line %d has an unknown colour (frame %d, val %d)", n+1, fN+1, vN+1)
				panic(msg)

			}

		}

	}

	return id
}
