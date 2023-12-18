package challenge

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
)

const STARTX = 1
const STARTY = 1

type LandType int

const (
	Undug LandType = iota
	Wall
	Centre
)

func (l LandType) String() string {
	switch l {
	case Wall:
		return "#"
	case Centre:
		return "x"

	default:
		return "."
	}
}

type DigOrder struct {
	Direction cartesian.Direction
	Length    int
}

func Challenge(input string) string {
	lines := strings.Split(input, "\n")

	grid := cartesian.CoordinateGrid[LandType]{}

	orders := []DigOrder{}
	var maxWest, maxNorth, curWest, curNorth int

	for _, l := range lines {
		if l == "" {
			continue
		}
		order := parseLine(l)

		switch order.Direction {
		case cartesian.North:
			curNorth += order.Length
		case cartesian.East:
			curWest -= order.Length
		case cartesian.South:
			curNorth -= order.Length
		case cartesian.West:
			curWest += order.Length

		}

		if curWest > maxWest {
			maxWest = curWest
		}

		if curNorth > maxNorth {
			maxNorth = curNorth
		}

		orders = append(orders, order)
	}

	curPos := cartesian.Coordinate{maxWest, maxNorth}

	totalNums := 0

	for _, order := range orders {
		totalNums += order.Length
	}

	for _, order := range orders {
		fmt.Println(order)
		for i := 0; i < order.Length; i++ {
			grid.Add(curPos, Wall)
			curPos = curPos.TransformInDirection(order.Direction)
		}
	}

	fmt.Println(&grid)

	grid.FloodFill(cartesian.Coordinate{STARTX, STARTY}, Wall, Centre)

	totalArea := 0

	for _, row := range grid.GetRows() {
		for _, dt := range row {
			switch dt {
			case Wall, Centre:
				totalArea++
			}
		}
	}

	return strconv.Itoa(totalArea)
}

func parseLine(l string) DigOrder {
	components := strings.Split(l, " ")

	stripped := components[2][2 : len(components[2])-1]
	var direction cartesian.Direction

	length, err := strconv.ParseInt(string(stripped[:5]), 16, 32)
	if err != nil {
		panic("int conv error")
	}

	switch stripped[5] {
	case '0':
		direction = cartesian.East
	case '1':
		direction = cartesian.South
	case '2':
		direction = cartesian.West
	case '3':
		direction = cartesian.North

	}

	return DigOrder{
		Direction: direction,
		Length:    int(length),
	}
}
