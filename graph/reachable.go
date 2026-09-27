package graph

import "slices"

func (g *Graph) Path(from, to int, mask Rel) []int {
	prev := map[int]int{}
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
				prev[to] = n
				path := []int{to}
				for path[len(path)-1] != from {
					path = append(path, prev[path[len(path)-1]])
				}
				slices.Reverse(path)
				return path
			}
			if !seen[next] {
				seen[next] = true
				prev[next] = n
				queue = append(queue, next)
			}
		}
	}
	return nil
}

func (g *Graph) Reachable(from, to int, mask Rel) bool {
	return g.Path(from, to, mask) != nil
}
