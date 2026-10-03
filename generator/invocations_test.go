package generator

import (
	"testing"

	"elle-go/history"
)

func TestInvocationsUniqueAppends(t *testing.T) {
	cfg := Config{Keys: []string{"x", "y", "z"}, Txns: 100, MaxMops: 4}
	ops := Invocations(cfg, 1)
	if len(ops) != cfg.Txns {
		t.Fatalf("got %d ops, want %d", len(ops), cfg.Txns)
	}
	seen := map[int]bool{}
	for _, op := range ops {
		if op.Type != history.Invoke || len(op.Mops) == 0 || len(op.Mops) > cfg.MaxMops {
			t.Fatalf("bad op: %+v", op)
		}
		for _, m := range op.Mops {
			if m.Type == history.Append {
				if seen[m.App] {
					t.Fatalf("duplicate append value %d", m.App)
				}
				seen[m.App] = true
			}
		}
	}
}
