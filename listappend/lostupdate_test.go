package listappend

import (
	"slices"
	"testing"

	"elle-go/history"
)

func TestLostUpdate(t *testing.T) {
	cases := []struct {
		name string
		h    history.History
		want bool
	}{
		{"same read then append", history.History{
			txn(rd("x"), app("x", 1)),
			txn(rd("x"), app("x", 2)),
		}, true},
		{"different reads", history.History{
			txn(rd("x"), app("x", 1)),
			txn(rd("x", 1), app("x", 2)),
		}, false},
		{"blind append", history.History{
			txn(rd("x"), app("x", 1)),
			txn(app("x", 2)),
		}, false},
		{"append before read", history.History{
			txn(rd("x"), app("x", 1)),
			txn(app("x", 2), rd("x", 2)),
		}, false},
		{"read only", history.History{
			txn(rd("x"), app("x", 1)),
			txn(rd("x")),
		}, false},
		{"different keys", history.History{
			txn(rd("x"), app("x", 1)),
			txn(rd("y"), app("y", 2)),
		}, false},
		{"second append same txn", history.History{
			txn(rd("x"), app("x", 1), app("x", 2)),
		}, false},
		{"first read counts", history.History{
			txn(rd("x"), app("x", 1)),
			txn(rd("x"), rd("x", 1), app("x", 2)),
		}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, got := LostUpdate(c.h); got != c.want {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestLostUpdateWitness(t *testing.T) {
	h := history.History{
		txn(app("x", 1)),
		txn(rd("x", 1), app("x", 2)),
		txn(rd("y")),
		txn(rd("x", 1), app("x", 3)),
	}
	c, ok := LostUpdate(h)
	if !ok || c.Key != "x" || !slices.Equal(c.Value, []int{1}) || !slices.Equal(c.Ops, []int{1, 3}) {
		t.Fatalf("got %+v, %v", c, ok)
	}
}

func TestLostUpdateIgnoresFailed(t *testing.T) {
	failed := txn(rd("x"), app("x", 2))
	failed.Type = history.Fail
	if _, ok := LostUpdate(history.History{txn(rd("x"), app("x", 1)), failed}); ok {
		t.Fatal("failed op should not count")
	}
}
