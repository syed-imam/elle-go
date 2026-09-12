package listappend

import "elle-go/history"

func VersionOrder(h history.History) map[string][]int {
	orders := map[string][]int{}
	for _, op := range h {
		if op.Type != history.Ok {
			continue
		}
		for _, mop := range op.Mops {
			if mop.Type != history.Read {
				continue
			}
			if len(mop.Read) > len(orders[mop.Key]) {
				orders[mop.Key] = mop.Read
			}
		}
	}
	return orders
}
