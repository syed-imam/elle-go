package listappend

import (
	"slices"
	"sort"

	"elle-go/history"
)

type OrderConflict struct {
	Key     string
	Shorter []int
	Longer  []int
}

func IncompatibleOrder(h history.History) (OrderConflict, bool) {
	reads := map[string][][]int{}
	for _, op := range h {
		if op.Type != history.Ok {
			continue
		}
		for _, mop := range op.Mops {
			if mop.Type == history.Read {
				reads[mop.Key] = append(reads[mop.Key], mop.Read)
			}
		}
	}
	keys := make([]string, 0, len(reads))
	for k := range reads {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		rs := reads[k]
		sort.SliceStable(rs, func(i, j int) bool { return len(rs[i]) < len(rs[j]) })
		for i := 0; i+1 < len(rs); i++ {
			if !slices.Equal(rs[i], rs[i+1][:len(rs[i])]) {
				return OrderConflict{Key: k, Shorter: rs[i], Longer: rs[i+1]}, true
			}
		}
	}
	return OrderConflict{}, false
}
