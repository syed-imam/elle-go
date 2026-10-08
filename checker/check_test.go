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
			v, err := Check(c.h)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if v.Anomaly != c.want {
				t.Errorf("got %v, want %v", v.Anomaly, c.want)
			}
		})
	}
}

func TestCheckRejectsUnwrittenRead(t *testing.T) {
	h := history.History{
		{Process: 1, Type: history.Ok, Mops: []history.Mop{{Type: history.Read, Key: "x", Read: []int{99}}}},
	}
	if _, err := Check(h); err == nil {
		t.Fatal("expected error for read of unwritten value")
	}
}

func TestCheckIncompatibleOrder(t *testing.T) {
	h := history.History{
		{Process: 1, Type: history.Ok, Mops: []history.Mop{{Type: history.Append, Key: "x", App: 1}}},
		{Process: 2, Type: history.Ok, Mops: []history.Mop{{Type: history.Append, Key: "x", App: 2}}},
		{Process: 3, Type: history.Ok, Mops: []history.Mop{{Type: history.Read, Key: "x", Read: []int{1, 2}}}},
		{Process: 4, Type: history.Ok, Mops: []history.Mop{{Type: history.Read, Key: "x", Read: []int{2, 1}}}},
	}
	v, err := Check(h)
	if err != nil {
		t.Fatal(err)
	}
	if v.Anomaly != anomaly.IncompatibleOrder || v.Key != "x" {
		t.Fatalf("got %v on key %q", v.Anomaly, v.Key)
	}
}
