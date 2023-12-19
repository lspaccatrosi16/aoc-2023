package challenge

import (
	"strconv"
	"strings"

	"github.com/lspaccatrosi16/go-libs/structures/stack"
)

type PartStatus int

const (
	Processing PartStatus = iota
	Accepted
	Rejected
)

func (p PartStatus) String() string {
	switch p {
	case Processing:
		return "Processing"
	case Accepted:
		return "Accepted"
	case Rejected:
		return "Rejected"
	default:
		return "Invalid"
	}
}

type Part struct {
	Status PartStatus
	Values map[string]int
}

type Workflow struct {
	Name       string
	Conditions []WorkflowCondition
}

func (w *Workflow) Eval(p *Part) string {
	for _, cond := range w.Conditions {
		valDesired := cond.DesiredVal
		valHad := p.Values[cond.ValName]
		var satisfied bool
		switch cond.Op {
		case GT:
			satisfied = valHad > valDesired
		case LT:
			satisfied = valHad < valDesired
		case NOOP:
			satisfied = true
		}

		if satisfied {
			return cond.Destination
		}
	}
	return w.Conditions[len(w.Conditions)-1].Destination
}

type CmpOp int

const (
	NOOP CmpOp = iota
	GT
	LT
)

func (c CmpOp) String() string {
	switch c {
	case GT:
		return ">"
	case LT:
		return "<"
	default:
		return "x"
	}
}

type WorkflowCondition struct {
	ValName     string
	DesiredVal  int
	Op          CmpOp
	Destination string
}

type Worker struct {
	Workflows map[string]Workflow
}

type PartOrder struct {
	Part         *Part
	WorkflowName string
}

func (w *Worker) EvalPart(part *Part) {
	stack := stack.NewStack[PartOrder]()

	stack.Push(PartOrder{
		Part:         part,
		WorkflowName: "in",
	})

	for {
		v, ok := stack.Pop()
		if !ok {
			break
		}

		workflow, ok := w.Workflows[v.WorkflowName]
		if !ok {
			panic("workflow not found: " + v.WorkflowName)
		}
		nextDest := workflow.Eval(v.Part)

		switch nextDest {
		case "R":
			v.Part.Status = Rejected
		case "A":
			v.Part.Status = Accepted
		default:
			stack.Push(PartOrder{
				Part:         v.Part,
				WorkflowName: nextDest,
			})
		}

	}

}

func Challenge(input string) string {
	sections := strings.Split(input, "\n\n")

	workflows := map[string]Workflow{}

	for _, l := range strings.Split(sections[0], "\n") {
		if l == "" {
			continue
		}
		wf := parseWorkflow(l)
		workflows[wf.Name] = wf
	}

	totalParts := 0
	worker := Worker{
		Workflows: workflows,
	}

	for m := 1; m <= 4000; m++ {
		for a := 1; a <= 4000; a++ {
			for s := 1; s <= 4000; s++ {
				for x := 1; x <= 4000; x++ {
					part := Part{
						Status: Processing,
						Values: map[string]int{
							"x": x,
							"s": s,
							"a": a,
							"m": m,
						},
					}
					worker.EvalPart(&part)
					if part.Status == Accepted {
						totalParts++
					}
				}
			}
		}
	}

	return strconv.Itoa(totalParts)
}

func parseWorkflow(s string) Workflow {
	components := strings.Split(s, "{")
	name := components[0]
	stripped := components[1][:len(components[1])-1]

	ruleStrs := strings.Split(stripped, ",")
	rules := []WorkflowCondition{}

	for _, str := range ruleStrs {
		rule := WorkflowCondition{}
		items := strings.Split(str, ":")
		if len(items) > 1 {
			rule.Destination = items[1]
			rule.ValName = string(items[0][0])
			switch items[0][1] {
			case '>':
				rule.Op = GT
			case '<':
				rule.Op = LT
			}

			vInt, err := strconv.Atoi(items[0][2:])
			if err != nil {
				panic("int conv error")
			}
			rule.DesiredVal = vInt
		} else {
			rule.Op = NOOP
			rule.Destination = items[0]
		}
		rules = append(rules, rule)
	}

	workflow := Workflow{
		Name:       name,
		Conditions: rules,
	}
	return workflow
}
