package challenge

import (
	"fmt"
	"strconv"
	"strings"

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

func (n *Node) HandlePulse(idFrom string, state bool) ([]string, bool) {
	// fmt.Printf("%s -%t-> %s\n", idFrom, state, n.id)

	switch n.nt {
	case Broadcaster:
		return n.destinations, state
	case FlipFlop:
		if !state {
			n.state["default"] = !n.state["default"]
			return n.destinations, n.state["default"]
		}
		return []string{}, false
	case Conjunction:
		n.state[idFrom] = state
		for _, v := range n.state {
			if !v {
				return n.destinations, true
			}
		}
		return n.destinations, false
	case Output:
		return []string{}, state
	}

	fmt.Println()

	panic("missed case")
}

type System struct {
	Cycles       int
	CycleResults []Cycle
	Cache        map[Cycle]int
	Nodes        map[string]*Node
}

type NodeOrder struct {
	Id    string
	From  string
	State bool
}

func (s *System) Run(n int) {
	for i := 0; i < n; i++ {
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
		// if v, ok := s.Cache[thisCycle]; ok {
		// 	// fmt.Println("cache hit found", v)
		// }
		s.Cache[thisCycle] = i
	}
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
}

func Challenge(input string) string {
	lines := strings.Split(input, "\n")

	conjs := map[string]bool{}

	nodes := map[string]*Node{}
	for _, l := range lines {
		if l == "" {
			continue
		}

		node := parseLine(l)
		if node.nt == Conjunction {
			conjs[node.id] = true
		}

		nodes[node.id] = node
	}

	nodes["output"] = &Node{
		nt: Output,
		id: "output",
	}

	conjipts := map[string][]string{}
	for _, n := range nodes {
		for _, dst := range n.destinations {
			if conjs[dst] {
				conjipts[dst] = append(conjipts[dst], n.id)
			}
		}
	}

	for k, v := range conjipts {
		for _, i := range v {
			nodes[k].state[i] = false
		}
	}

	system := System{
		Nodes: nodes,
		Cache: map[Cycle]int{},
	}

	system.Run(1000)

	high, low := system.TallyPresses()

	fmt.Println(high, low)

	total := high * low

	return strconv.Itoa(total)
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
