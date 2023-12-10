package challenge

import (
	"fmt"
	"strconv"
	"strings"
)

type Coordinate [2]int

func (c Coordinate) String() string {
	return fmt.Sprintf("(%d, %d)", c[0], c[1])
}

func (c Coordinate) IsSame(o Coordinate) bool {
	return c[0] == o[0] && c[1] == o[1]
}

type Pipe struct {
	Position  Coordinate
	Connected [2]*Pipe
	Type      rune
}

func (p *Pipe) PointsTo(pipe *Pipe) bool {
	return p.Connected[0] == pipe || p.Connected[1] == pipe
}

func (p *Pipe) PointsToCoord(c Coordinate) bool {
	return p.Connected[0].Position == c || p.Connected[1].Position == c
}

type TraverseData struct {
	LastNode *Pipe
	CurNode  *Pipe
}

type Traverser struct {
	InitialNode *Pipe
	TravData    TraverseData
}

func (t *Traverser) Traverse() *map[int]map[int]bool {
	t.TravData = TraverseData{
		LastNode: t.InitialNode,
		CurNode:  t.InitialNode.Connected[0],
	}

	pointMap := map[int]map[int]bool{}

	ctr := 1
	for {
		addTo2DMap(&pointMap, true, t.TravData.CurNode.Position[0], t.TravData.CurNode.Position[1])
		for i := 0; i < 2; i++ {
			if n := t.TravData.CurNode.Connected[i]; n != t.TravData.CurNode && n != t.TravData.LastNode {
				t.TravData.LastNode = t.TravData.CurNode
				t.TravData.CurNode = n
				break
			}
		}

		if t.TravData.CurNode.Position.IsSame(t.InitialNode.Position) {
			break
		}
		ctr++
	}

	addTo2DMap(&pointMap, true, t.TravData.CurNode.Position[0], t.TravData.CurNode.Position[1])
	return &pointMap
}

func Challenge(input string) string {
	iptLines := strings.Split(input, "\n")

	maxY := len(iptLines)
	maxX := len(iptLines[0])

	pipeMap := map[int]map[int]*Pipe{}
	pipeList := []*Pipe{}

	possibleNestMap := map[int]map[int]bool{}

	var startingPipe *Pipe

	for y, l := range iptLines {
		for x, c := range l {
			pos := Coordinate{x, y}
			if c == '.' {
				addTo2DMap(&possibleNestMap, true, x, y)
				continue
			}
			pipe := &Pipe{
				Position: pos,
				Type:     c,
			}

			if c == 'S' {
				startingPipe = pipe
			}
			addTo2DMap(&pipeMap, pipe, x, y)
			pipeList = append(pipeList, pipe)
		}
	}

	for i := 0; i < len(pipeList); i++ {
		connectedCoords := [2]Coordinate{}
		switch pipeList[i].Type {
		case '|':
			connectedCoords[0] = transformCoordinate(pipeList[i].Position, 0, -1)
			connectedCoords[1] = transformCoordinate(pipeList[i].Position, 0, 1)
		case '-':
			connectedCoords[0] = transformCoordinate(pipeList[i].Position, -1, 0)
			connectedCoords[1] = transformCoordinate(pipeList[i].Position, 1, 0)
		case 'L':
			connectedCoords[0] = transformCoordinate(pipeList[i].Position, 0, -1)
			connectedCoords[1] = transformCoordinate(pipeList[i].Position, 1, 0)
		case 'J':
			connectedCoords[0] = transformCoordinate(pipeList[i].Position, 0, -1)
			connectedCoords[1] = transformCoordinate(pipeList[i].Position, -1, 0)
		case '7':
			connectedCoords[0] = transformCoordinate(pipeList[i].Position, -1, 0)
			connectedCoords[1] = transformCoordinate(pipeList[i].Position, 0, 1)
		case 'F':
			connectedCoords[0] = transformCoordinate(pipeList[i].Position, 1, 0)
			connectedCoords[1] = transformCoordinate(pipeList[i].Position, 0, 1)
		}

		for j := 0; j < 2; j++ {
			pipeList[i].Connected[j] = retrieveFromMap(&pipeMap, connectedCoords[j])
		}
	}

	pointsToStart := []*Pipe{}

	if p := retrieveFromMap(&pipeMap, transformCoordinate(startingPipe.Position, 0, -1)); p != nil && p.PointsTo(startingPipe) {
		pointsToStart = append(pointsToStart, p)
	}

	if p := retrieveFromMap(&pipeMap, transformCoordinate(startingPipe.Position, 0, 1)); p != nil && p.PointsTo(startingPipe) {
		pointsToStart = append(pointsToStart, p)
	}

	if p := retrieveFromMap(&pipeMap, transformCoordinate(startingPipe.Position, -1, 0)); p != nil && p.PointsTo(startingPipe) {
		pointsToStart = append(pointsToStart, p)
	}

	if p := retrieveFromMap(&pipeMap, transformCoordinate(startingPipe.Position, 1, 0)); p != nil && p.PointsTo(startingPipe) {
		pointsToStart = append(pointsToStart, p)
	}

	if len(pointsToStart) != 2 {
		panic("points to start has unexpected length")
	}

	startingPipe.Connected[0] = pointsToStart[0]
	startingPipe.Connected[1] = pointsToStart[1]

	traverser := Traverser{
		InitialNode: startingPipe,
	}

	loopCoords := traverser.Traverse()

	newL := make([][]rune, maxY)
	numCells := 0

	for y := 0; y < maxY; y++ {
		l := make([]rune, maxX)
		walls := 0
		for x := 0; x < maxX; x++ {
			if retrieveFromMap(loopCoords, Coordinate{x, y}) {
				l[x] = '#'
				switch iptLines[y][x] {
				case '|', 'S':
					walls++
				case 'J':
					walls++
				case 'L':
					walls++
				}
			} else if walls%2 == 1 {
				numCells++
				l[x] = 'I'
			} else {
				l[x] = 'O'
			}
		}
		newL[y] = l
	}

	reprStr := ""

	for i := 0; i < len(newL); i++ {
		reprStr += string(newL[i]) + "\n"
	}

	fmt.Println(reprStr)

	return strconv.Itoa(numCells)
}

func addTo2DMap[T any](m *map[int]map[int]T, p T, x, y int) {
	if val, exists := (*m)[x]; exists {
		val[y] = p
		(*m)[x] = val
	} else {
		newY := map[int]T{
			y: p,
		}
		(*m)[x] = newY
	}
}

func transformCoordinate(c Coordinate, xoffset, yoffset int) Coordinate {
	return Coordinate{c[0] + xoffset, c[1] + yoffset}
}

func retrieveFromMap[T any](m *map[int]map[int]T, c Coordinate) T {
	return (*m)[c[0]][c[1]]
}
