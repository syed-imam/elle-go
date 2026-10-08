package listappend

import (
	"slices"
	"testing"

	"elle-go/history"
)

func readOp(key string, vals ...int) history.Op {
	return history.Op{Type: history.Ok, Mops: []history.Mop{{Type: history.Read, Key: key, Read: vals}}}
}

func TestIncompatibleOrder(t *testing.T) {
	cases := []struct {
		name string
		h    history.History
		want bool
	}{
		{"prefixes", history.History{readOp("x", 1), readOp("x", 1, 2, 3), readOp("x"), readOp("x", 1, 2)}, false},
		{"diverging", history.History{readOp("x", 1, 2), readOp("x", 1, 3, 2)}, true},
		{"same length differ", history.History{readOp("x", 1, 2), readOp("x", 2, 1)}, true},
		{"other key fine", history.History{readOp("x", 1), readOp("y", 2)}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, got := IncompatibleOrder(c.h); got != c.want {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestIncompatibleOrderWitness(t *testing.T) {
	c, ok := IncompatibleOrder(history.History{readOp("x", 1, 3, 2), readOp("x", 1, 2)})
	if !ok || c.Key != "x" || !slices.Equal(c.Shorter, []int{1, 2}) || !slices.Equal(c.Longer, []int{1, 3, 2}) {
		t.Fatalf("got %+v, %v", c, ok)
	}
}

func TestIncompatibleOrderIgnoresFailed(t *testing.T) {
	failed := readOp("x", 2, 1)
	failed.Type = history.Fail
	if _, ok := IncompatibleOrder(history.History{readOp("x", 1, 2), failed}); ok {
		t.Fatal("failed op should not count")
	}
}
