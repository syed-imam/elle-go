package graph

import (
	"testing"

	"elle-go/fixtures"
	"elle-go/listappend"
)

func fromEdges(edges []listappend.Edge) *Graph {
	g := New()
	for _, e := range edges {
		g.AddEdge(e.From, e.To)
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
