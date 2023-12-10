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
	NumSteps int
}

type Traverser struct {
	InitialNode *Pipe
	TravData    TraverseData
}

func (t *Traverser) TraverseUntilCollide() int {
	t.TravData = TraverseData{
		LastNode: t.InitialNode,
		CurNode:  t.InitialNode.Connected[0],
		NumSteps: 1,
	}
	for {

		// total steps to get back to start / 2

		for i := 0; i < 2; i++ {
			if n := t.TravData.CurNode.Connected[i]; n != t.TravData.CurNode && n != t.TravData.LastNode {
				t.TravData.LastNode = t.TravData.CurNode
				t.TravData.CurNode = n
				t.TravData.NumSteps++
				break
			}
		}

		if t.TravData.CurNode.Position.IsSame(t.InitialNode.Position) {
			break
		}
	}

	return t.TravData.NumSteps / 2
}

func Challenge(input string) string {
	lines := strings.Split(input, "\n")

	pipeMap := map[int]map[int]*Pipe{}
	pipeList := []*Pipe{}

	var startingPipe *Pipe

	for y, l := range lines {
		for x, c := range l {
			pos := Coordinate{x, y}
			if c == '.' {
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
			pipeList[i].Connected[j] = retrieveFromCoordinate(&pipeMap, connectedCoords[j])
		}
	}

	pointsToStart := []*Pipe{}

	if p := retrieveFromCoordinate(&pipeMap, transformCoordinate(startingPipe.Position, 0, -1)); p != nil && p.PointsTo(startingPipe) {
		pointsToStart = append(pointsToStart, p)
	}

	if p := retrieveFromCoordinate(&pipeMap, transformCoordinate(startingPipe.Position, 0, 1)); p != nil && p.PointsTo(startingPipe) {
		pointsToStart = append(pointsToStart, p)
	}

	if p := retrieveFromCoordinate(&pipeMap, transformCoordinate(startingPipe.Position, -1, 0)); p != nil && p.PointsTo(startingPipe) {
		pointsToStart = append(pointsToStart, p)
	}

	if p := retrieveFromCoordinate(&pipeMap, transformCoordinate(startingPipe.Position, 1, 0)); p != nil && p.PointsTo(startingPipe) {
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

	furthest := traverser.TraverseUntilCollide()

	return strconv.Itoa(furthest)
}

func addTo2DMap(m *map[int]map[int]*Pipe, p *Pipe, x, y int) {
	if val, exists := (*m)[x]; exists {
		val[y] = p
		(*m)[x] = val
	} else {
		newY := map[int]*Pipe{
			y: p,
		}
		(*m)[x] = newY
	}
}

func transformCoordinate(c Coordinate, xoffset, yoffset int) Coordinate {
	return Coordinate{c[0] + xoffset, c[1] + yoffset}
}

func retrieveFromCoordinate(m *map[int]map[int]*Pipe, c Coordinate) *Pipe {
	return (*m)[c[0]][c[1]]
}
