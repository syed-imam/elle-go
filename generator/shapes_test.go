package generator

import (
	"testing"

	"elle-go/history"
)

func TestFromShapes(t *testing.T) {
	shapes := []history.Op{
		{Mops: []history.Mop{{Type: history.Read, Key: "y"}, {Type: history.Append, Key: "x"}}},
		{Mops: []history.Mop{{Type: history.Read, Key: "x"}, {Type: history.Append, Key: "y"}}},
	}
	ops := FromShapes(shapes, 50, 1)
	if len(ops) != 50 {
		t.Fatalf("got %d ops, want 50", len(ops))
	}
	seen := map[int]bool{}
	for _, op := range ops {
		if op.Type != history.Invoke || len(op.Mops) != 2 || op.Mops[0].Type != history.Read || op.Mops[1].Type != history.Append {
			t.Fatalf("op does not match a shape: %+v", op)
		}
		if seen[op.Mops[1].App] {
			t.Fatalf("duplicate append value %d", op.Mops[1].App)
		}
		seen[op.Mops[1].App] = true
	}
	if shapes[0].Mops[1].App != 0 {
		t.Fatal("FromShapes mutated its input")
	}
}
