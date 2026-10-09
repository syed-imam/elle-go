package listappend

import (
	"slices"

	"github.com/syed-imam/elle-go/history"
)

type InternalCase struct {
	Op       int
	Key      string
	Expected []int
	Prefixed bool
	Read     []int
}

func Internal(h history.History) (InternalCase, bool) {
	for i, op := range h {
		if op.Type != history.Ok {
			continue
		}
		if c, ok := opInternal(op); ok {
			c.Op = i
			return c, true
		}
	}
	return InternalCase{}, false
}

type expectation struct {
	vals     []int
	prefixed bool
}

func opInternal(op history.Op) (InternalCase, bool) {
	state := map[string]expectation{}
	for _, mop := range op.Mops {
		s, known := state[mop.Key]
		switch mop.Type {
		case history.Append:
			if !known {
				s.prefixed = true
			}
			s.vals = append(slices.Clone(s.vals), mop.App)
			state[mop.Key] = s
		case history.Read:
			if known && !matches(s, mop.Read) {
				return InternalCase{Key: mop.Key, Expected: s.vals, Prefixed: s.prefixed, Read: mop.Read}, true
			}
			state[mop.Key] = expectation{vals: mop.Read}
		}
	}
	return InternalCase{}, false
}

func matches(s expectation, read []int) bool {
	if !s.prefixed {
		return slices.Equal(s.vals, read)
	}
	i := len(read) - len(s.vals)
	return i >= 0 && slices.Equal(s.vals, read[i:])
}
