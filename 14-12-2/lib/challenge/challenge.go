package challenge

import (
	"fmt"
	"strconv"
	"strings"
)

const SnapshotInterval = 1_000_000
const CycleNumber = 1_000_000_000

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

func (cg *CoordGrid) arrs() [][]rune {
	if cg == nil {
		return nil
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

	return lines
}

func (cg *CoordGrid) String() string {
	lines := cg.arrs()

	outStr := ""
	for _, l := range lines {
		outStr += string(l) + "\n"
	}
	return outStr
}

func (cg *CoordGrid) ToHash() string {
	arrs := cg.arrs()

	hashStr := ""

	for _, l := range arrs {
		for _, c := range l {
			hashStr += string(c)
		}
	}

	return hashStr

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

	totalWeight := 0

	for i := 0; i < CycleNumber; i++ {
		weight, done := runTiltCycle(movableRocks, &coordGrid, i+1, maxY)
		totalWeight += weight
		if done {
			totalWeight = weight
			break
		}
	}

	fmt.Println(coordGrid.String())

	return strconv.Itoa(totalWeight)
}

var cache = make(map[string]int)
var arrangements = make(map[int]int)

func runTiltCycle(mr []*Rock, cg *CoordGrid, n, maxY int) (int, bool) {
	combinations := [4][2]int{
		{0, -1}, {-1, 0}, {0, 1}, {1, 0},
	}

	for _, combo := range combinations {
		for {
			changes := 0
			for i := 0; i < len(mr); i++ {
				newPos := mr[i].Position.Transform(combo[0], combo[1])
				curEnq := cg.Get(newPos)

				if curEnq != nil && curEnq.Type == Empty {
					cg.Swap(mr[i].Position, newPos)
					curEnq.Position = mr[i].Position
					mr[i].Position = newPos
					changes++
				}
			}
			if changes == 0 {
				break
			}
		}
	}

	weight := 0

	for _, r := range mr {
		weight += maxY - r.Position[1]
	}

	if v, ok := cache[cg.ToHash()]; ok {
		start := (CycleNumber)%(n-v) + 2*(n-v)
		return arrangements[start], true
	}

	cache[cg.ToHash()] = n
	arrangements[n] = weight

	return weight, false

}

// 108, 109

// 93742
