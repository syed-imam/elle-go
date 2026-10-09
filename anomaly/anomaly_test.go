package anomaly

import (
	"testing"

	"github.com/syed-imam/elle-go/fixtures"
	"github.com/syed-imam/elle-go/graph"
	"github.com/syed-imam/elle-go/listappend"
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

func isClosedCycle(g *graph.Graph, c []int) bool {
	if len(c) < 2 {
		return false
	}
	for i := range c {
		if g.EdgeRel(c[i], c[(i+1)%len(c)]) == 0 {
			return false
		}
	}
	return true
}

func TestCheckWitnessCycle(t *testing.T) {
	cases := []struct {
		name    string
		build   func(*graph.Graph)
		anomaly Anomaly
	}{
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
		v := Check(g)
		if v.Anomaly != c.anomaly {
			t.Errorf("%s: Anomaly = %v, want %v", c.name, v.Anomaly, c.anomaly)
		}
		if v.Level != Requires(c.anomaly) {
			t.Errorf("%s: Level = %v, want %v", c.name, v.Level, Requires(c.anomaly))
		}
		if !isClosedCycle(g, v.Cycle) {
			t.Errorf("%s: Cycle %v is not a closed loop", c.name, v.Cycle)
		}
	}

	g := graph.New()
	g.AddEdge(1, 2, graph.WW)
	if v := Check(g); v.Anomaly != None || v.Cycle != nil {
		t.Errorf("acyclic: got %v cycle %v, want None nil", v.Anomaly, v.Cycle)
	}
}

func TestVerdictString(t *testing.T) {
	if s := (Verdict{Anomaly: None}).String(); s != "no anomaly found" {
		t.Errorf("None.String() = %q, want %q", s, "no anomaly found")
	}

	v := Verdict{Anomaly: G2, Level: Serializable, Cycle: []int{0, 1}}
	want := "G2: transactions 0 → 1 → 0 (requires serializable or stronger)"
	if got := v.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}

	v = IncompatibleOrderVerdict("x", []int{1, 2}, []int{2, 1})
	want = `incompatible-order: key "x" read as [1 2] and [2 1], which no single order of appends explains (requires read committed or stronger)`
	if got := v.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}

	v = InternalVerdict(3, "x", []int{1}, false, []int{1, 2})
	want = `internal: transaction 3 read key "x" as [1 2], but its own earlier reads and appends imply [1] (requires snapshot isolation or stronger)`
	if got := v.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}

	v = LostUpdateVerdict("x", []int{}, []int{1, 2})
	want = `lost-update: transactions 1 and 2 both read key "x" as [], then appended to it (requires snapshot isolation or stronger)`
	if got := v.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}

	v = InternalVerdict(0, "x", []int{2}, true, []int{1})
	want = `internal: transaction 0 read key "x" as [1], but its own earlier reads and appends imply …[2] (requires snapshot isolation or stronger)`
	if got := v.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
