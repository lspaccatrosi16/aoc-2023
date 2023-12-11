package challenge

import (
	"math"
	"strconv"
	"strings"
)

const ExpandFactor = 1_000_000

const NumWorkers = 8

type Coordinate [2]int

func (c Coordinate) Transform(x, y int) Coordinate {
	return Coordinate{c[0] + x, c[1] + y}
}

func (c Coordinate) DistanceStepwise(other Coordinate) int {
	xDiff := math.Abs(float64(other[0] - c[0]))
	yDiff := math.Abs(float64(other[1] - c[1]))

	return int(xDiff + yDiff)
}

type Galaxy struct {
	Position Coordinate
}

type Universe []*Galaxy

type GalaxyPair [2]*Galaxy

func Challenge(input string) string {
	lines := strings.Split(input, "\n")

	maxY := len(lines)
	maxX := len(lines[0])

	universeInX := map[int]bool{}
	universeInY := map[int]bool{}

	for y, l := range lines {
		for x, c := range l {
			if c == '#' {
				universeInX[x] = true
				universeInY[y] = true
			}
		}
	}

	universe := Universe{}

	curY := 0
	for y := 0; y < maxY; y++ {
		curX := 0
		for x := 0; x < maxX; x++ {
			if !universeInX[x] {
				curX += ExpandFactor
			} else {
				if lines[y][x] == '#' {
					galaxy := Galaxy{
						Position: Coordinate{curX, curY},
					}
					universe = append(universe, &galaxy)
				}
				curX++
			}
		}
		if !universeInY[y] {
			curY += ExpandFactor
		} else {
			curY++
		}
	}

	tally := 0

	for i := 0; i < len(universe); i++ {
		for j := i + 1; j < len(universe); j++ {
			tally += universe[i].Position.DistanceStepwise(universe[j].Position)
		}
	}

	return strconv.Itoa(tally)
}
