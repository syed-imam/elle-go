package fixtures

import (
	"testing"

	"elle-go/history"
)

func TestFixtures(t *testing.T) {
	cases := []struct {
		name    string
		build   func() history.History
		wantOps int
	}{
		{"serializable", Serializable, 3},
		{"write-skew", WriteSkew, 2},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := len(c.build()); got != c.wantOps {
				t.Errorf("got %d ops, want %d", got, c.wantOps)
			}
		})
	}
}
