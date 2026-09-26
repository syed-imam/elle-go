package anomaly

import "elle-go/graph"

type Anomaly int

const (
	None Anomaly = iota
	G0
	G1c
	GSingle
	G2
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
	default:
		return "none"
	}
}

func Classify(g *graph.Graph) Anomaly {
	if g.Filter(graph.WW).FindCycle() != nil {
		return G0
	}
	if closesLoop(g, graph.WR) {
		return G1c
	}
	if closesLoop(g, graph.RW) {
		return GSingle
	}
	if g.FindCycle() != nil {
		return G2
	}
	return None
}

func closesLoop(g *graph.Graph, via graph.Rel) bool {
	for _, e := range g.EdgesWith(via) {
		if g.Reachable(e[1], e[0], graph.WW|graph.WR) {
			return true
		}
	}
	return false
}
