package graph

type Rel uint8

const (
	WW Rel = 1 << iota
	WR
	RW
)

type Graph struct {
	adj  map[int][]int
	rels map[[2]int]Rel
}

func New() *Graph {
	return &Graph{adj: map[int][]int{}, rels: map[[2]int]Rel{}}
}

func (g *Graph) AddEdge(from, to int, rel Rel) {
	g.adj[from] = append(g.adj[from], to)
	g.rels[[2]int{from, to}] |= rel
}

func (g *Graph) EdgeRel(from, to int) Rel {
	return g.rels[[2]int{from, to}]
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
			r.AddEdge(to, from, g.EdgeRel(from, to))
		}
	}
	return r
}
