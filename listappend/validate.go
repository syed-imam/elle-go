package listappend

import (
	"fmt"

	"elle-go/history"
)

func Validate(h history.History) error {
	writers := appendedBy(h)
	for i, op := range h {
		if op.Type != history.Ok {
			continue
		}
		for _, mop := range op.Mops {
			if mop.Type != history.Read {
				continue
			}
			for _, v := range mop.Read {
				if _, ok := writers[mop.Key][v]; !ok {
					return fmt.Errorf("op %d read %v from key %q, but no ok or info op appended it", i, v, mop.Key)
				}
			}
		}
	}
	return nil
}
