package graph

type Graph struct {
	adj map[int][]int
}

func New() *Graph {
	return &Graph{adj: map[int][]int{}}
}

func (g *Graph) AddEdge(from, to int) {
	g.adj[from] = append(g.adj[from], to)
}

func (g *Graph) nodes() []int {
	set := map[int]bool{}
	for from, tos := range g.adj {
		set[from] = true
		for _, to := range tos {
			set[to] = true
		}
	}
	out := make([]int, 0, len(set))
	for n := range set {
		out = append(out, n)
	}
	return out
}

func (g *Graph) reverse() *Graph {
	r := New()
	for from, tos := range g.adj {
		for _, to := range tos {
			r.AddEdge(to, from)
		}
	}
	return r
}
