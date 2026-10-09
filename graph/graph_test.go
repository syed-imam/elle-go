package graph

import (
	"testing"

	"github.com/syed-imam/elle-go/fixtures"
	"github.com/syed-imam/elle-go/listappend"
)

func rel(t listappend.EdgeType) Rel {
	switch t {
	case listappend.WR:
		return WR
	case listappend.RW:
		return RW
	default:
		return WW
	}
}

func fromEdges(edges []listappend.Edge) *Graph {
	g := New()
	for _, e := range edges {
		g.AddEdge(e.From, e.To, rel(e.Type))
	}
	return g
}

func TestHasCycle(t *testing.T) {
	if !fromEdges(listappend.Dependencies(fixtures.WriteSkew())).HasCycle() {
		t.Error("write-skew: expected a cycle, found none")
	}
	if fromEdges(listappend.Dependencies(fixtures.Serializable())).HasCycle() {
		t.Error("serializable: expected no cycle, found one")
	}
}
