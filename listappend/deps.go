package listappend

import "github.com/syed-imam/elle-go/history"

func Dependencies(h history.History) []Edge {
	writers := appendedBy(h)
	orders := VersionOrder(h)
	edges := []Edge{}
	observed := observedValues(orders)

	for key, order := range orders {
		for i := 0; i+1 < len(order); i++ {
			from := writers[key][order[i]]
			to := writers[key][order[i+1]]
			if from != to {
				edges = append(edges, Edge{From: from, To: to, Type: WW})
			}
		}
	}

	for i, op := range h {
		if op.Type != history.Ok {
			continue
		}
		for _, mop := range op.Mops {
			if mop.Type != history.Read {
				continue
			}
			seen := map[int]bool{}
			for _, v := range mop.Read {
				seen[v] = true
				if w := writers[mop.Key][v]; w != i {
					edges = append(edges, Edge{From: w, To: i, Type: WR})
				}
			}
			for v, w := range writers[mop.Key] {
				if !seen[v] && w != i && (h[w].Type == history.Ok || observed[mop.Key][v]) {
					edges = append(edges, Edge{From: i, To: w, Type: RW})
				}
			}
		}
	}

	return edges
}

func appendedBy(h history.History) map[string]map[int]int {
	writers := map[string]map[int]int{}
	for i, op := range h {
		if op.Type != history.Ok && op.Type != history.Info {
			continue
		}
		for _, mop := range op.Mops {
			if mop.Type != history.Append {
				continue
			}
			if writers[mop.Key] == nil {
				writers[mop.Key] = map[int]int{}
			}
			writers[mop.Key][mop.App] = i
		}
	}
	return writers
}

func observedValues(orders map[string][]int) map[string]map[int]bool {
	observed := map[string]map[int]bool{}
	for key, order := range orders {
		observed[key] = map[int]bool{}
		for _, v := range order {
			observed[key][v] = true
		}
	}
	return observed
}
