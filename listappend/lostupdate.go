package listappend

import (
	"fmt"

	"github.com/syed-imam/elle-go/history"
)

type LostUpdateCase struct {
	Key   string
	Value []int
	Ops   []int
}

type keyValue struct {
	key   string
	value string
}

func LostUpdate(h history.History) (LostUpdateCase, bool) {
	first := map[keyValue]int{}
	for i, op := range h {
		if op.Type != history.Ok {
			continue
		}
		reads := map[string][]int{}
		written := map[string]bool{}
		for _, mop := range op.Mops {
			read, seen := reads[mop.Key]
			switch {
			case mop.Type == history.Read && !seen:
				reads[mop.Key] = mop.Read
			case mop.Type == history.Append && seen && !written[mop.Key]:
				written[mop.Key] = true
				kv := keyValue{mop.Key, fmt.Sprint(read)}
				if j, ok := first[kv]; ok {
					return LostUpdateCase{Key: mop.Key, Value: read, Ops: []int{j, i}}, true
				}
				first[kv] = i
			}
		}
	}
	return LostUpdateCase{}, false
}
