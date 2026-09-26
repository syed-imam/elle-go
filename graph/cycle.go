package graph

import "slices"

func (g *Graph) FindCycle() []int {
	state := map[int]int{}
	var path []int
	var cycle []int

	var visit func(n int) bool
	visit = func(n int) bool {
		state[n] = 1
		path = append(path, n)
		for _, next := range g.adj[n] {
			if state[next] == 1 {
				cycle = slices.Clone(path[slices.Index(path, next):])
				return true
			}
			if state[next] == 0 && visit(next) {
				return true
			}
		}
		path = path[:len(path)-1]
		state[n] = 2
		return false
	}

	for _, n := range g.nodes() {
		if state[n] == 0 && visit(n) {
			return cycle
		}
	}
	return nil
}
