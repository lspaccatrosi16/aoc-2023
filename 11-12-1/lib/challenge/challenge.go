package challenge

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
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

	expandedLines := []string{}

	for y := 0; y < maxY; y++ {
		thisLine := ""
		for x := 0; x < maxX; x++ {
			if !universeInX[x] {
				thisLine += strings.Repeat(".", ExpandFactor)
			} else {
				thisLine += string(lines[y][x])
			}
		}
		if !universeInY[y] {
			for i := 0; i < ExpandFactor; i++ {
				expandedLines = append(expandedLines, thisLine)
			}
		} else {
			expandedLines = append(expandedLines, thisLine)
		}
	}

	wg := &sync.WaitGroup{}
	results := make(chan *Galaxy)

	fmt.Printf("Expanded Lines: %d\n", len(expandedLines))

	for y, l := range expandedLines {
		wg.Add(1)
		ls := lineScan{l, y}
		go scanLine(wg, ls, results)
	}

	go monitor(wg, results)

	tally := 0
	universe := Universe{}

	for g := range results {
		universe = append(universe, g)
	}

	for i := 0; i < len(universe); i++ {
		for j := i + 1; j < len(universe); j++ {
			tally += universe[i].Position.DistanceStepwise(universe[j].Position)
		}
	}

	return strconv.Itoa(tally)
}

func monitor(wg *sync.WaitGroup, results chan *Galaxy) {
	wg.Wait()
	close(results)
}

type lineScan struct {
	l string
	y int
}

func scanLine(wg *sync.WaitGroup, ls lineScan, results chan *Galaxy) {
	defer wg.Done()
	if ls.y%100_000 == 0 {
		fmt.Printf("Work Line %d\n", ls.y/100_000)
	}
	for x, c := range ls.l {
		if c == '#' {
			galaxy := Galaxy{Position: Coordinate{x, ls.y}}
			results <- &galaxy
		}
	}
}
