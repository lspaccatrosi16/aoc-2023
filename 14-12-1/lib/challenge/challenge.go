package challenge

import (
	"fmt"
	"strconv"
	"strings"
)

type Coordinate [2]int

func (c Coordinate) Transform(xoff, yoff int) Coordinate {
	return Coordinate{c[0] + xoff, c[1] + yoff}
}

type RockType int

const (
	Round RockType = iota
	Square
	Empty
)

func (r RockType) String() string {
	switch r {
	case Round:
		return "O"
	case Square:
		return "#"
	default:
		return "."
	}
}

type Rock struct {
	Type     RockType
	Position Coordinate
}

type CoordGrid map[int]map[int]*Rock

func (cg *CoordGrid) Add(r *Rock, c Coordinate) {
	if cg == nil {
		cg = new(CoordGrid)
	}

	if v, ok := (*cg)[c[0]]; ok {
		v[c[1]] = r
	} else {
		(*cg)[c[0]] = map[int]*Rock{c[1]: r}
	}
}

func (c *CoordGrid) Swap(c1, c2 Coordinate) {
	if c == nil {
		return
	}
	(*c)[c1[0]][c1[1]], (*c)[c2[0]][c2[1]] = (*c)[c2[0]][c2[1]], (*c)[c1[0]][c1[1]]
}

func (cg *CoordGrid) Get(c Coordinate) *Rock {
	if cg == nil {
		return nil
	}
	return (*cg)[c[0]][c[1]]
}

func (cg *CoordGrid) String() string {
	if cg == nil {
		return ""
	}

	maxX := 0
	maxY := 0

	for x, l := range *cg {
		if x > maxX {
			maxX = x
		}
		for y := range l {
			if y > maxY {
				maxY = y
			}
		}
	}

	lines := make([][]rune, maxY+1)

	for i := 0; i < len(lines); i++ {
		lines[i] = make([]rune, maxX+1)
	}

	for x, l := range *cg {
		for y, c := range l {
			lines[y][x] = rune(c.Type.String()[0])
		}
	}

	outStr := ""
	for _, l := range lines {
		outStr += string(l) + "\n"
	}
	return outStr
}

func Challenge(input string) string {
	coordGrid := CoordGrid{}

	lines := strings.Split(input, "\n")
	maxY := len(lines)

	movableRocks := []*Rock{}

	for y, l := range lines {
		for x, c := range l {
			var rt RockType
			if c == 'O' {
				rt = Round
			} else if c == '#' {
				rt = Square
			} else {
				rt = Empty
			}

			rock := &Rock{
				Type:     rt,
				Position: Coordinate{x, y},
			}
			coordGrid.Add(rock, rock.Position)
			if rock.Type == Round {
				movableRocks = append(movableRocks, rock)
			}
		}
	}

	for {
		changes := 0

		for i := 0; i < len(movableRocks); i++ {
			newPos := movableRocks[i].Position.Transform(0, -1)
			if newPos[1] < 0 {
				continue
			}
			curEnq := coordGrid.Get(newPos)

			if curEnq != nil && curEnq.Type == Empty {
				coordGrid.Swap(movableRocks[i].Position, newPos)
				curEnq.Position = movableRocks[i].Position
				movableRocks[i].Position = newPos
				changes++
			}
		}

		if changes == 0 {
			break
		}
	}

	fmt.Println(coordGrid.String())

	totalWeight := 0

	for _, r := range movableRocks {
		totalWeight += maxY - r.Position[1]
	}

	return strconv.Itoa(totalWeight)
}
