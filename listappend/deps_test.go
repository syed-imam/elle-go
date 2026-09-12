package listappend

import (
	"testing"

	"elle-go/fixtures"
)

func has(edges []Edge, from, to int, t EdgeType) bool {
	for _, e := range edges {
		if e.From == from && e.To == to && e.Type == t {
			return true
		}
	}
	return false
}

func countType(edges []Edge, t EdgeType) int {
	n := 0
	for _, e := range edges {
		if e.Type == t {
			n++
		}
	}
	return n
}

func TestWriteSkewHasCycle(t *testing.T) {
	edges := Dependencies(fixtures.WriteSkew())
	if !has(edges, 0, 1, RW) || !has(edges, 1, 0, RW) {
		t.Errorf("expected rw cycle 0<->1, got %v", edges)
	}
}

func TestSerializableNoRW(t *testing.T) {
	edges := Dependencies(fixtures.Serializable())
	if n := countType(edges, RW); n != 0 {
		t.Errorf("expected no rw edges, got %d in %v", n, edges)
	}
	if !has(edges, 0, 1, WW) {
		t.Errorf("expected ww edge 0->1, got %v", edges)
	}
}
