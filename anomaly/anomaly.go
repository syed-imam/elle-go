package anomaly

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"elle-go/graph"
)

type Anomaly int

const (
	None Anomaly = iota
	G0
	G1c
	GSingle
	G2
	IncompatibleOrder
)

func (a Anomaly) String() string {
	switch a {
	case G0:
		return "G0"
	case G1c:
		return "G1c"
	case GSingle:
		return "G-single"
	case G2:
		return "G2"
	case IncompatibleOrder:
		return "incompatible-order"
	default:
		return "none"
	}
}

type Verdict struct {
	Anomaly Anomaly
	Level   Level
	Cycle   []int
	Key     string
	Reads   [][]int
}

func IncompatibleOrderVerdict(key string, a, b []int) Verdict {
	return Verdict{Anomaly: IncompatibleOrder, Level: Requires(IncompatibleOrder), Key: key, Reads: [][]int{a, b}}
}

func Classify(g *graph.Graph) Anomaly {
	return Check(g).Anomaly
}

func Check(g *graph.Graph) Verdict {
	if c := g.Filter(graph.WW).FindCycle(); c != nil {
		return newVerdict(G0, c)
	}
	if c := closingCycle(g, graph.WR); c != nil {
		return newVerdict(G1c, c)
	}
	if c := closingCycle(g, graph.RW); c != nil {
		return newVerdict(GSingle, c)
	}
	if c := g.FindCycle(); c != nil {
		return newVerdict(G2, c)
	}
	return newVerdict(None, nil)
}

func newVerdict(a Anomaly, cycle []int) Verdict {
	return Verdict{Anomaly: a, Level: Requires(a), Cycle: cycle}
}

func closingCycle(g *graph.Graph, via graph.Rel) []int {
	for _, e := range g.EdgesWith(via) {
		if p := g.Path(e[1], e[0], graph.WW|graph.WR); p != nil {
			return append([]int{e[0]}, p[:len(p)-1]...)
		}
	}
	return nil
}

func (v Verdict) String() string {
	if v.Anomaly == None {
		return "no anomaly found"
	}
	if v.Anomaly == IncompatibleOrder {
		return fmt.Sprintf("%v: key %q read as %v and %v, which no single order of appends explains (requires %v or stronger)",
			v.Anomaly, v.Key, v.Reads[0], v.Reads[1], v.Level)
	}
	loop := append(slices.Clone(v.Cycle), v.Cycle[0])
	nodes := make([]string, len(loop))
	for i, n := range loop {
		nodes[i] = strconv.Itoa(n)
	}
	return fmt.Sprintf("%v: transactions %s (requires %v or stronger)",
		v.Anomaly, strings.Join(nodes, " → "), v.Level)
}
