package graph

import (
	"slices"
	"testing"
)

func TestNodes(t *testing.T) {
	g := New()
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)

	got := g.nodes()
	slices.Sort(got)
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("nodes = %v, want [1 2 3]", got)
	}
}

func TestReverse(t *testing.T) {
	g := New()
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)

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
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)
	g.AddEdge(3, 1)
	g.AddEdge(3, 4)
	g.AddEdge(4, 5)

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
	acyclic.AddEdge(1, 2)
	acyclic.AddEdge(2, 3)
	if acyclic.HasCycle() {
		t.Error("acyclic graph reported a cycle")
	}

	cyclic := New()
	cyclic.AddEdge(1, 2)
	cyclic.AddEdge(2, 1)
	if !cyclic.HasCycle() {
		t.Error("cyclic graph reported no cycle")
	}
}

func TestFinishOrder(t *testing.T) {
	g := New()
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)
	g.AddEdge(1, 3)

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
