package checker

import (
	"elle-go/anomaly"
	"elle-go/graph"
	"elle-go/history"
	"elle-go/listappend"
)

func Check(h history.History) anomaly.Verdict {
	g := graph.New()
	for _, e := range listappend.Dependencies(h) {
		g.AddEdge(e.From, e.To, toRel(e.Type))
	}
	return anomaly.Check(g)
}

func toRel(t listappend.EdgeType) graph.Rel {
	switch t {
	case listappend.WR:
		return graph.WR
	case listappend.RW:
		return graph.RW
	default:
		return graph.WW
	}
}
