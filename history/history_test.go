package history

import (
	"reflect"
	"testing"
)

func TestKeys(t *testing.T) {
	h := History{
		{Mops: []Mop{{Type: Read, Key: "y"}, {Type: Append, Key: "x"}}},
		{Mops: []Mop{{Type: Append, Key: "y"}, {Type: Read, Key: "z"}}},
	}
	if got := h.Keys(); !reflect.DeepEqual(got, []string{"x", "y", "z"}) {
		t.Fatalf("Keys() = %v", got)
	}
}
