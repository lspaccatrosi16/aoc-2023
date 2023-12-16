package challenge

import (
	"strconv"
	"strings"

	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
)

type TileType int

const (
	Unset TileType = iota
	Empty
	AngleRL
	AngleLR
	SplitHorizontal
	SplitVertical
)

func (t TileType) String() string {
	switch t {
	case AngleRL:
		return "/"
	case AngleLR:
		return "\\"
	case SplitHorizontal:
		return "-"
	case SplitVertical:
		return "|"
	case Empty:
		return "."
	default:
		return "X"
	}
}

func (t TileType) NewDirection(currentDir BeamDirection) []BeamDirection {
	if t == Empty {
		return []BeamDirection{currentDir}
	}

	if t == SplitHorizontal && (currentDir == East || currentDir == West) {
		return []BeamDirection{currentDir}
	}

	if t == SplitVertical && (currentDir == North || currentDir == South) {
		return []BeamDirection{currentDir}
	}

	if t == SplitHorizontal {
		return []BeamDirection{East, West}
	}

	if t == SplitVertical {
		return []BeamDirection{North, South}
	}

	if t == AngleRL {
		switch currentDir {
		case North:
			return []BeamDirection{East}
		case East:
			return []BeamDirection{North}
		case South:
			return []BeamDirection{West}
		case West:
			return []BeamDirection{South}
		}
	}

	if t == AngleLR {
		switch currentDir {
		case North:
			return []BeamDirection{West}
		case East:
			return []BeamDirection{South}
		case South:
			return []BeamDirection{East}
		case West:
			return []BeamDirection{North}
		}
	}

	return []BeamDirection{}
}

type BeamDirection int

const (
	North BeamDirection = iota
	East
	South
	West
)

func (b BeamDirection) String() string {
	switch b {
	case North:
		return "North"
	case East:
		return "East"
	case South:
		return "South"
	case West:
		return "West"
	default:
		return "Invalid"
	}
}

type Beam struct {
	VisitedPositions map[cartesian.Coordinate]bool
	CurrentPos       cartesian.Coordinate
	CurrentDirection BeamDirection
	Children         []*Beam
	Grid             *cartesian.CoordinateGrid[TileType]
}

type visitedPoint struct {
	Direction BeamDirection
	Position  cartesian.Coordinate
}

var vm map[visitedPoint]bool = map[visitedPoint]bool{}

func (b *Beam) Navigate(initialDirection BeamDirection, startCoordinate cartesian.Coordinate) {
	b.CurrentPos = startCoordinate
	initialPoint := b.Grid.Get(startCoordinate)

	startingDirection := initialPoint.NewDirection(initialDirection)

	switch len(startingDirection) {
	case 0:
		return
	case 1:
		b.CurrentDirection = startingDirection[0]
	case 2:
		beam0 := &Beam{Grid: b.Grid}
		beam1 := &Beam{Grid: b.Grid}

		b.Children = append(b.Children, beam0, beam1)
		beam0.Navigate(startingDirection[0], startCoordinate)
		beam1.Navigate(startingDirection[1], startCoordinate)
		return
	}

	b.CurrentDirection = startingDirection[0]

	b.VisitedPositions = map[cartesian.Coordinate]bool{}
	b.Children = []*Beam{}

	for {
		vpData := visitedPoint{
			Direction: b.CurrentDirection,
			Position:  b.CurrentPos,
		}

		if v, exists := vm[vpData]; exists && v {
			return
		}

		vm[vpData] = true

		var trX, trY int

		switch b.CurrentDirection {
		case North:
			trX, trY = 0, -1
		case East:
			trX, trY = 1, 0
		case South:
			trX, trY = 0, 1
		case West:
			trX, trY = -1, 0
		}

		newPos := b.CurrentPos.Transform(trX, trY)
		nextPoint := b.Grid.Get(newPos)

		newDirections := nextPoint.NewDirection(b.CurrentDirection)

		b.VisitedPositions[b.CurrentPos] = true

		switch len(newDirections) {
		case 0:
			return
		case 1:
			b.CurrentDirection = newDirections[0]
			b.CurrentPos = newPos
		case 2:
			beam0 := &Beam{Grid: b.Grid}
			beam1 := &Beam{Grid: b.Grid}

			b.Children = append(b.Children, beam0, beam1)
			beam0.Navigate(newDirections[0], newPos)
			beam1.Navigate(newDirections[1], newPos)
			return
		}
	}
}

func (b *Beam) EnergisedSquares() map[cartesian.Coordinate]bool {
	m := b.VisitedPositions

	for _, c := range b.Children {
		cm := c.EnergisedSquares()

		for k, v := range cm {
			if v {
				m[k] = v
			}
		}
	}

	return m
}

func Challenge(input string) string {
	lines := strings.Split(input, "\n")
	grid := &cartesian.CoordinateGrid[TileType]{}

	for y, l := range lines {
		for x, c := range l {
			var tt TileType
			switch c {
			case '/':
				tt = AngleRL
			case '\\':
				tt = AngleLR
			case '-':
				tt = SplitHorizontal
			case '|':
				tt = SplitVertical
			default:
				tt = Empty
			}

			grid.Add(cartesian.Coordinate{x, y}, tt)
		}
	}
	maxEnergised := 0
	masterBeam := Beam{Grid: grid}

	masterBeam.Navigate(East, cartesian.Coordinate{0, 0})
	maxEnergised = len(masterBeam.EnergisedSquares())

	return strconv.Itoa(maxEnergised)
}
