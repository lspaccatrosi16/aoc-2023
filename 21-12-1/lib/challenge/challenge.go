package challenge

import (
	"strconv"
	"strings"

	"github.com/lspaccatrosi16/go-libs/algorithms/graph"
	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
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

	bfsRes, err := graph.RunBfs(startNode, g, StepNum)
	if err != nil {
		panic(err)
	}

	possibleIds := []string{}

	for k, d := range bfsRes.Dist {
		if d%2 == 0 {
			possibleIds = append(possibleIds, k)
		}
	}

	return strconv.Itoa(len(possibleIds))
}
