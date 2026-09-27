package graph

import (
	"slices"
	"testing"
)

func TestNodes(t *testing.T) {
	g := New()
	g.AddEdge(1, 2, WW)
	g.AddEdge(2, 3, WW)

	got := g.nodes()
	slices.Sort(got)
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("nodes = %v, want [1 2 3]", got)
	}
}

func TestReverse(t *testing.T) {
	g := New()
	g.AddEdge(1, 2, WW)
	g.AddEdge(2, 3, WW)

	r := g.reverse()
	if !slices.Equal(r.adj[2], []int{1}) {
		t.Errorf("reverse adj[2] = %v, want [1]", r.adj[2])
	}
	if !slices.Equal(r.adj[3], []int{2}) {
		t.Errorf("reverse adj[3] = %v, want [2]", r.adj[3])
	}
}

func TestSCCs(t *testing.T) {
	g := New()
	g.AddEdge(1, 2, WW)
	g.AddEdge(2, 3, WW)
	g.AddEdge(3, 1, WW)
	g.AddEdge(3, 4, WW)
	g.AddEdge(4, 5, WW)

	comps := g.SCCs()
	if len(comps) != 3 {
		t.Fatalf("got %d components, want 3: %v", len(comps), comps)
	}

	var big []int
	for _, c := range comps {
		if len(c) == 3 {
			big = slices.Clone(c)
		}
	}
	slices.Sort(big)
	if !slices.Equal(big, []int{1, 2, 3}) {
		t.Errorf("largest SCC = %v, want [1 2 3]", big)
	}
}

func TestHasCycleUnit(t *testing.T) {
	acyclic := New()
	acyclic.AddEdge(1, 2, WW)
	acyclic.AddEdge(2, 3, WW)
	if acyclic.HasCycle() {
		t.Error("acyclic graph reported a cycle")
	}

	cyclic := New()
	cyclic.AddEdge(1, 2, WW)
	cyclic.AddEdge(2, 1, WW)
	if !cyclic.HasCycle() {
		t.Error("cyclic graph reported no cycle")
	}
}

func TestFindCycle(t *testing.T) {
	acyclic := New()
	acyclic.AddEdge(1, 2, WW)
	acyclic.AddEdge(2, 3, WW)
	if c := acyclic.FindCycle(); c != nil {
		t.Errorf("acyclic graph returned cycle %v", c)
	}

	cyclic := New()
	cyclic.AddEdge(1, 2, WW)
	cyclic.AddEdge(2, 3, WW)
	cyclic.AddEdge(3, 1, WW)
	cyclic.AddEdge(3, 4, WW)

	got := cyclic.FindCycle()
	if len(got) < 2 {
		t.Fatalf("cycle = %v, want a closed loop", got)
	}
	for i := range got {
		next := got[(i+1)%len(got)]
		if !slices.Contains(cyclic.adj[got[i]], next) {
			t.Errorf("cycle %v: no edge %d->%d", got, got[i], next)
		}
	}
}

func TestEdgeRel(t *testing.T) {
	g := New()
	g.AddEdge(1, 2, WW)
	g.AddEdge(1, 2, WR)

	got := g.EdgeRel(1, 2)
	if got&WW == 0 || got&WR == 0 {
		t.Errorf("EdgeRel(1,2) = %b, want WW and WR set", got)
	}
	if got&RW != 0 {
		t.Errorf("EdgeRel(1,2) = %b, should not have RW", got)
	}
	if g.EdgeRel(2, 1) != 0 {
		t.Errorf("EdgeRel(2,1) = %b, want empty", g.EdgeRel(2, 1))
	}
}

func TestFilter(t *testing.T) {
	g := New()
	g.AddEdge(1, 2, WR)
	g.AddEdge(2, 1, WW)

	if g.Filter(WW).FindCycle() != nil {
		t.Error("Filter(WW): no ww-only cycle should exist here")
	}
	if g.Filter(WW|WR).FindCycle() == nil {
		t.Error("Filter(WW|WR): a ww/wr cycle should exist")
	}

	pure := New()
	pure.AddEdge(1, 2, WW)
	pure.AddEdge(2, 1, WW)
	if pure.Filter(WW).FindCycle() == nil {
		t.Error("Filter(WW): ww-only cycle should be found")
	}
}

func TestReachable(t *testing.T) {
	g := New()
	g.AddEdge(1, 2, WW)
	g.AddEdge(2, 3, WR)
	g.AddEdge(3, 4, RW)

	if !g.Reachable(1, 3, WW|WR) {
		t.Error("1 should reach 3 via ww/wr")
	}
	if g.Reachable(1, 4, WW|WR) {
		t.Error("1 should not reach 4 without crossing the rw edge")
	}
	if !g.Reachable(1, 4, WW|WR|RW) {
		t.Error("1 should reach 4 when rw is allowed")
	}

	c := New()
	c.AddEdge(1, 2, WW)
	c.AddEdge(2, 1, WR)
	if !c.Reachable(1, 1, WW|WR) {
		t.Error("1 should reach itself around a ww/wr cycle")
	}
}

func TestEdgesWith(t *testing.T) {
	g := New()
	g.AddEdge(1, 2, WW)
	g.AddEdge(2, 3, RW)
	g.AddEdge(3, 1, RW)
	g.AddEdge(3, 1, WR)

	rw := g.EdgesWith(RW)
	if len(rw) != 2 {
		t.Fatalf("EdgesWith(RW) = %v, want 2 edges", rw)
	}
	if !slices.Contains(rw, [2]int{2, 3}) || !slices.Contains(rw, [2]int{3, 1}) {
		t.Errorf("EdgesWith(RW) = %v, want [2 3] and [3 1]", rw)
	}

	ww := g.EdgesWith(WW)
	if len(ww) != 1 || ww[0] != [2]int{1, 2} {
		t.Errorf("EdgesWith(WW) = %v, want [[1 2]]", ww)
	}
}

func TestPath(t *testing.T) {
	g := New()
	g.AddEdge(1, 2, WW)
	g.AddEdge(2, 3, WR)
	g.AddEdge(3, 4, RW)

	if p := g.Path(1, 3, WW|WR); !slices.Equal(p, []int{1, 2, 3}) {
		t.Errorf("Path(1,3,ww|wr) = %v, want [1 2 3]", p)
	}
	if p := g.Path(1, 4, WW|WR); p != nil {
		t.Errorf("Path(1,4,ww|wr) = %v, want nil (rw edge excluded)", p)
	}
	if p := g.Path(1, 4, WW|WR|RW); !slices.Equal(p, []int{1, 2, 3, 4}) {
		t.Errorf("Path(1,4,all) = %v, want [1 2 3 4]", p)
	}
}

func TestFinishOrder(t *testing.T) {
	g := New()
	g.AddEdge(1, 2, WW)
	g.AddEdge(2, 3, WW)
	g.AddEdge(1, 3, WW)

	pos := map[int]int{}
	for i, n := range g.finishOrder() {
		pos[n] = i
	}
	for from, tos := range g.adj {
		for _, to := range tos {
			if pos[from] <= pos[to] {
				t.Errorf("edge %d->%d: %d should finish after %d", from, to, from, to)
			}
		}
	}
}
