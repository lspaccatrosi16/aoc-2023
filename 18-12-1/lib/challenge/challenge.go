package challenge

import (
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
	Colour    [3]uint8
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
	for _, order := range orders {
		for i := 0; i < order.Length; i++ {
			grid.Add(curPos, Wall)
			curPos = curPos.TransformInDirection(order.Direction)
		}
	}
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

	var direction cartesian.Direction
	switch components[0] {
	case "U":
		direction = cartesian.North
	case "D":
		direction = cartesian.South
	case "L":
		direction = cartesian.West
	case "R":
		direction = cartesian.East
	}

	length, err := strconv.Atoi(components[1])

	if err != nil {
		panic("int conv error")
	}

	stripped := components[2][2 : len(components[2])-1]

	if len(stripped) != 6 {
		panic("stripped wrong length")
	}

	colours := [3]uint8{}

	for i := 0; i < 3; i++ {
		colStr := string(stripped[2*i : (2*i)+1])
		col, err := strconv.ParseUint(colStr, 16, 8)
		if err != nil {
			panic("uint parse error")
		}
		colours[i] = uint8(col)
	}

	return DigOrder{
		Colour:    colours,
		Direction: direction,
		Length:    length,
	}
}
