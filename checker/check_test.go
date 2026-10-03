package checker

import (
	"testing"

	"elle-go/anomaly"
	"elle-go/fixtures"
	"elle-go/history"
)

func TestCheck(t *testing.T) {
	cases := []struct {
		name string
		h    history.History
		want anomaly.Anomaly
	}{
		{"serializable", fixtures.Serializable(), anomaly.None},
		{"write skew", fixtures.WriteSkew(), anomaly.G2},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Check(c.h).Anomaly; got != c.want {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}
