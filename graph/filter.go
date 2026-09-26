package graph

func (g *Graph) Filter(mask Rel) *Graph {
	f := New()
	for from, tos := range g.adj {
		for _, to := range tos {
			if r := g.EdgeRel(from, to) & mask; r != 0 {
				f.AddEdge(from, to, r)
			}
		}
	}
	return f
}
