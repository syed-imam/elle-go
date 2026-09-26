package anomaly

import (
	"testing"

	"elle-go/fixtures"
	"elle-go/graph"
	"elle-go/listappend"
)

func TestClassifyShapes(t *testing.T) {
	cases := []struct {
		name  string
		build func(*graph.Graph)
		want  Anomaly
	}{
		{"none", func(g *graph.Graph) {
			g.AddEdge(1, 2, graph.WW)
		}, None},
		{"g0", func(g *graph.Graph) {
			g.AddEdge(1, 2, graph.WW)
			g.AddEdge(2, 1, graph.WW)
		}, G0},
		{"g1c", func(g *graph.Graph) {
			g.AddEdge(1, 2, graph.WR)
			g.AddEdge(2, 1, graph.WW)
		}, G1c},
		{"gsingle", func(g *graph.Graph) {
			g.AddEdge(1, 2, graph.RW)
			g.AddEdge(2, 1, graph.WW)
		}, GSingle},
		{"g2", func(g *graph.Graph) {
			g.AddEdge(1, 2, graph.RW)
			g.AddEdge(2, 1, graph.RW)
		}, G2},
	}
	for _, c := range cases {
		g := graph.New()
		c.build(g)
		if got := Classify(g); got != c.want {
			t.Errorf("%s: Classify = %v, want %v", c.name, got, c.want)
		}
	}
}

func build(edges []listappend.Edge) *graph.Graph {
	g := graph.New()
	for _, e := range edges {
		var r graph.Rel
		switch e.Type {
		case listappend.WR:
			r = graph.WR
		case listappend.RW:
			r = graph.RW
		default:
			r = graph.WW
		}
		g.AddEdge(e.From, e.To, r)
	}
	return g
}

func TestClassifyFixtures(t *testing.T) {
	if got := Classify(build(listappend.Dependencies(fixtures.WriteSkew()))); got != G2 {
		t.Errorf("write-skew: Classify = %v, want G2", got)
	}
	if got := Classify(build(listappend.Dependencies(fixtures.Serializable()))); got != None {
		t.Errorf("serializable: Classify = %v, want None", got)
	}
}
