package challenge

import (
	"fmt"
	"strings"

	ga "github.com/lspaccatrosi16/go-libs/algorithms/graph"
	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
	gs "github.com/lspaccatrosi16/go-libs/structures/graph"
)

const StepNum = 64

type LandType int

const (
	Rock LandType = iota
	Plot
)

func Challenge(input string) string {
	lines := strings.Split(input, "\n")

	grid := cartesian.CoordinateGrid[LandType]{}

	var start cartesian.Coordinate

	for y, l := range lines {
		for x, c := range l {
			var p LandType
			switch c {
			case 'S':
				p = Plot
				start = cartesian.Coordinate{x, y}
			case '.':
				p = Plot
			case '#':
				p = Rock
			}

			grid.Add(cartesian.Coordinate{x, y}, p)
		}
	}

	g, nm := grid.CreateGraph(false, []LandType{Plot})
	startNode := (*nm)[start]

	completeFill := findPossibleSteps(g, startNode, -1)

	n := 0
	for {
		n++
		tryFill := findPossibleSteps(g, startNode, n)
		if completeFill == tryFill {
			break
		}
	}

	fmt.Println(n)

	return ""
}

func findPossibleSteps(g *gs.Graph, startNode gs.GraphNode, n int) int {
	bfsRes, err := ga.RunBfs(startNode, g, n)
	if err != nil {
		panic(err)
	}

	spots := 0

	for _, d := range bfsRes.Dist {
		if d%2 == 0 {
			spots++
		}
	}
	return spots
}
