package graph

func (g *Graph) finishOrder() []int {
	seen := map[int]bool{}
	order := []int{}

	var visit func(n int)
	visit = func(n int) {
		seen[n] = true
		for _, next := range g.adj[n] {
			if !seen[next] {
				visit(next)
			}
		}
		order = append(order, n)
	}

	for _, n := range g.nodes() {
		if !seen[n] {
			visit(n)
		}
	}
	return order
}

func (g *Graph) SCCs() [][]int {
	order := g.finishOrder()
	rev := g.reverse()
	seen := map[int]bool{}
	comps := [][]int{}

	var collect func(n int, comp *[]int)
	collect = func(n int, comp *[]int) {
		seen[n] = true
		*comp = append(*comp, n)
		for _, next := range rev.adj[n] {
			if !seen[next] {
				collect(next, comp)
			}
		}
	}

	for i := len(order) - 1; i >= 0; i-- {
		n := order[i]
		if !seen[n] {
			comp := []int{}
			collect(n, &comp)
			comps = append(comps, comp)
		}
	}
	return comps
}

func (g *Graph) HasCycle() bool {
	for _, comp := range g.SCCs() {
		if len(comp) > 1 {
			return true
		}
	}
	return false
}
