package graph

func (g *Graph) Reachable(from, to int, mask Rel) bool {
	seen := map[int]bool{from: true}
	queue := []int{from}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		for _, next := range g.adj[n] {
			if g.EdgeRel(n, next)&mask == 0 {
				continue
			}
			if next == to {
				return true
			}
			if !seen[next] {
				seen[next] = true
				queue = append(queue, next)
			}
		}
	}
	return false
}
