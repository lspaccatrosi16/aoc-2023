package challenge

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lspaccatrosi16/go-libs/algorithms/maths"
	"github.com/lspaccatrosi16/go-libs/structures/mpq"
)

type NodeType int

const (
	Broadcaster NodeType = iota
	FlipFlop
	Conjunction
	Output
)

type Node struct {
	state        map[string]bool
	nt           NodeType
	id           string
	destinations []string
}

func (n *Node) GetState() bool {
	switch n.nt {
	case FlipFlop:
		return n.state["default"]
	case Conjunction:
		if n.id == "tg" {
			fmt.Println(n.state)
		}
		for _, v := range n.state {
			if !v {
				return true
			}
		}
		return false
	default:
		return false
	}

}

func (n *Node) HandlePulse(idFrom string, state bool) ([]string, bool) {
	switch n.nt {
	case Broadcaster:
		return n.destinations, state
	case FlipFlop:
		if !state {
			n.state["default"] = !n.state["default"]
			return n.destinations, n.GetState()
		}
		return []string{}, false
	case Conjunction:
		n.state[idFrom] = state
		return n.destinations, n.GetState()
	case Output:
		return []string{}, state
	}

	panic("missed case")
}

type NodeList map[string]*Node

func (n *NodeList) ResetState() {
	conjs := map[string]bool{}

	for k, v := range *n {
		if v.nt == Conjunction {
			conjs[k] = true
		}

		v.state = map[string]bool{}
	}

	conjipts := map[string][]string{}
	for _, node := range *n {
		for _, dst := range node.destinations {
			if conjs[dst] {
				conjipts[dst] = append(conjipts[dst], node.id)
			}
		}
	}

	for k, v := range conjipts {
		for _, i := range v {
			(*n)[k].state[i] = false
		}
	}
}

type System struct {
	Cycles       int
	CycleResults []Cycle
	Nodes        NodeList
	PrevList     map[string][]string
}

type NodeOrder struct {
	Id    string
	From  string
	State bool
}

func (s *System) AnalyseConjunction(output string) []string {
	leading := s.PrevList[output]

	dependents := []string{}

	for _, l := range leading {
		if s.Nodes[l].nt != Conjunction {
			return []string{output}
		}
		dependents = append(dependents, s.AnalyseConjunction(l)...)
	}

	return dependents
}

func (s *System) Run(names []string) []int {
	nameIdx := map[string]int{}

	for i, n := range names {
		nameIdx[n] = i
	}

	loops := make([]int, len(names))

	for {
		thisCycle := Cycle{}

		queue := mpq.Queue[NodeOrder]{}

		initialOrder := NodeOrder{
			Id:    "broadcaster",
			From:  "button",
			State: false,
		}
		queue.Add(initialOrder, 1)

		var outputState bool

		for queue.Len() != 0 {
			order := queue.Pop()

			if order.State {
				thisCycle.HighPresses++
			} else {
				thisCycle.LowPresses++
			}

			node, ok := s.Nodes[order.Id]
			if !ok {
				continue
			}
			destinations, evalState := node.HandlePulse(order.From, order.State)

			if node.nt == Output {
				outputState = evalState
			}

			for _, dst := range destinations {
				newOrder := NodeOrder{
					Id:    dst,
					State: evalState,
					From:  order.Id,
				}
				queue.Add(newOrder, 1)
			}
		}

		thisCycle.OutputState = outputState

		s.Cycles++
		s.CycleResults = append(s.CycleResults, thisCycle)

		set := true

		for k, v := range nameIdx {
			n := s.Nodes[k]
			if n.GetState() {
				fmt.Println("found occurance")
				loops[v] = s.Cycles
			}
		}

		for _, l := range loops {
			if l == 0 {
				set = false
			}
		}

		if set {
			break
		}
	}

	return loops
}

func (s *System) TallyPresses() (int, int) {
	totalHigh := 0
	totalLow := 0

	for _, c := range s.CycleResults {
		totalHigh += c.HighPresses
		totalLow += c.LowPresses
	}

	return totalHigh, totalLow
}

type Cycle struct {
	HighPresses int
	LowPresses  int
	OutputState bool
	RxPresses   int
}

func Challenge(input string) string {
	lines := strings.Split(input, "\n")

	nodes := NodeList{}
	for _, l := range lines {
		if l == "" {
			continue
		}

		node := parseLine(l)
		nodes[node.id] = node
	}

	nodes["output"] = &Node{
		nt: Output,
		id: "output",
	}

	nodes.ResetState()

	prev := map[string][]string{}

	for k, v := range nodes {
		for _, dst := range v.destinations {
			prev[dst] = append(prev[dst], k)
		}
	}

	system := System{
		Nodes:    nodes,
		PrevList: prev,
	}

	deps := system.AnalyseConjunction("rx")
	fmt.Println(deps)

	cycles := system.Run(deps)

	val := maths.Lcm(cycles...)

	return strconv.Itoa(val)
}

func parseLine(l string) *Node {
	components := strings.Split(l, " ")

	var nt NodeType
	var name string

	switch components[0][0] {
	case 'b':
		nt = Broadcaster
		name = "broadcaster"
	case '%':
		nt = FlipFlop
		name = components[0][1:]

	case '&':
		nt = Conjunction
		name = components[0][1:]
	}

	if components[1] != "->" {
		panic("expected connecter")
	}

	ending := strings.Join(components[2:], "")
	ids := strings.Split(ending, ",")

	node := &Node{
		state:        map[string]bool{},
		nt:           nt,
		id:           name,
		destinations: ids,
	}

	return node
}
