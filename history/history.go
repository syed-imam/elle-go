package history

import "sort"

type History []Op

func (h History) Keys() []string {
	seen := map[string]bool{}
	keys := []string{}
	for _, op := range h {
		for _, m := range op.Mops {
			if !seen[m.Key] {
				seen[m.Key] = true
				keys = append(keys, m.Key)
			}
		}
	}
	sort.Strings(keys)
	return keys
}
