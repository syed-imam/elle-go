package listappend

import (
	"slices"
	"testing"

	"elle-go/history"
)

func txn(mops ...history.Mop) history.Op {
	return history.Op{Type: history.Ok, Mops: mops}
}

func app(key string, v int) history.Mop {
	return history.Mop{Type: history.Append, Key: key, App: v}
}

func rd(key string, vals ...int) history.Mop {
	return history.Mop{Type: history.Read, Key: key, Read: vals}
}

func TestInternal(t *testing.T) {
	cases := []struct {
		name string
		op   history.Op
		want bool
	}{
		{"repeat read same", txn(rd("x", 1), rd("x", 1)), false},
		{"repeat read grew", txn(rd("x", 1), rd("x", 1, 2)), true},
		{"read own append", txn(rd("x", 1), app("x", 2), rd("x", 1, 2)), false},
		{"own append missing", txn(rd("x", 1), app("x", 2), rd("x", 1)), true},
		{"blind append then read suffix", txn(app("x", 2), rd("x", 1, 2)), false},
		{"blind append then read without it", txn(app("x", 2), rd("x", 1)), true},
		{"blind append then read too short", txn(app("x", 2), app("x", 3), rd("x", 3)), true},
		{"blind appends in order", txn(app("x", 2), app("x", 3), rd("x", 2, 3)), false},
		{"first read unconstrained", txn(rd("x", 5), rd("y", 1)), false},
		{"other key ignored", txn(rd("x", 1), app("y", 2), rd("x", 1)), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, got := Internal(history.History{c.op}); got != c.want {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestInternalWitness(t *testing.T) {
	h := history.History{txn(rd("x", 1)), txn(rd("x", 1), app("x", 2), rd("x", 1, 3, 2))}
	c, ok := Internal(h)
	if !ok || c.Op != 1 || c.Key != "x" || c.Prefixed || !slices.Equal(c.Expected, []int{1, 2}) || !slices.Equal(c.Read, []int{1, 3, 2}) {
		t.Fatalf("got %+v, %v", c, ok)
	}
}

func TestInternalIgnoresFailed(t *testing.T) {
	op := txn(rd("x", 1), rd("x", 2))
	op.Type = history.Fail
	if _, ok := Internal(history.History{op}); ok {
		t.Fatal("failed op should not count")
	}
}
