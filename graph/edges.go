package graph

func (g *Graph) EdgesWith(mask Rel) [][2]int {
	seen := map[[2]int]bool{}
	var out [][2]int
	for from, tos := range g.adj {
		for _, to := range tos {
			key := [2]int{from, to}
			if seen[key] || g.EdgeRel(from, to)&mask == 0 {
				continue
			}
			seen[key] = true
			out = append(out, key)
		}
	}
	return out
}
