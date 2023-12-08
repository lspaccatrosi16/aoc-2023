package challenge

import (
	"strconv"
	"strings"
)

func Challenge(input string) string {
	lines := strings.Split(input, "\n")

	instructionsLetters := lines[0]
	commands := lines[2:]

	instructions := []Instruction{}
	for i := 0; i < len(instructionsLetters); i++ {
		let := instructionsLetters[i]
		if let == ' ' || let == '\n' || let == '\t' || let == '\r' {
			continue
		}
		instr := new(Instruction).FromByte(let)
		instructions = append(instructions, instr)
	}

	protoNodeList := []*ProtoNode{}

	for _, cmd := range commands {
		comps := strings.Split(cmd, "=")
		id := strings.Trim(comps[0], " \n\r\t")
		pathways := []string{}
		curIdRead := ""

		for i := 0; i < len(comps[1]); i++ {
			c := comps[1][i]
			if c == ' ' || c == '(' {
				continue
			}

			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
				curIdRead += string(c)
			} else if c == ',' || c == ')' {
				pathways = append(pathways, curIdRead)
				curIdRead = ""
			}
		}

		pn := &ProtoNode{
			NodeId:   id,
			Pathways: pathways,
		}

		protoNodeList = append(protoNodeList, pn)
	}

	nodeMap := map[string]*Node{}

	for _, v := range protoNodeList {
		nodeMap[v.NodeId] = &Node{NodeId: v.NodeId}
	}

	for i := 0; i < len(protoNodeList); i++ {
		pn := protoNodeList[i]
		realNode := nodeMap[pn.NodeId]

		foundPathways := []*Node{}

		for _, pw := range pn.Pathways {
			n, found := nodeMap[pw]
			if found {
				foundPathways = append(foundPathways, n)
			} else {
				panic("Could not find pathway")
			}
		}

		realNode.Pathways = foundPathways
	}

	rootNode := nodeMap["AAA"]

	mapRt := Map{
		StartNode:    rootNode,
		Instructions: instructions,
	}

	mapRt.NavigateTo("ZZZ")

	return strconv.Itoa(mapRt.Steps)
}

type Instruction int

const (
	Left Instruction = iota
	Right
)

func (i Instruction) FromByte(b byte) Instruction {
	switch b {
	case 'L':
		return Left
	case 'R':
		return Right
	}

	return Instruction(-1)
}

func (i Instruction) String() string {
	switch i {
	case Left:
		return "Left"
	case Right:
		return "Right"
	}

	return "Invalid"
}

type Map struct {
	Steps        int
	StartNode    *Node
	Instructions []Instruction
}

func (m *Map) NavigateTo(id string) {
	currentNode := m.StartNode

	i := 0
	for {
		if currentNode.NodeId == id {
			break
		}

		if i >= len(m.Instructions) {
			i = 0
		}

		instr := m.Instructions[i]

		switch instr {
		case Left:
			currentNode = currentNode.Pathways[0]
		case Right:
			currentNode = currentNode.Pathways[1]
		}

		m.Steps++
		i++
	}
}

type Node struct {
	Pathways []*Node
	NodeId   string
}

type ProtoNode struct {
	Pathways []string
	NodeId   string
}
