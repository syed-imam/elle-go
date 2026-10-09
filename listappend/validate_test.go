package listappend

import (
	"testing"

	"github.com/syed-imam/elle-go/history"
)

func TestValidateUnknownRead(t *testing.T) {
	h := history.History{
		{Process: 1, Type: history.Ok, Mops: []history.Mop{{Type: history.Append, Key: "x", App: 1}}},
		{Process: 2, Type: history.Ok, Mops: []history.Mop{{Type: history.Read, Key: "y", Read: []int{99}}}},
	}
	if err := Validate(h); err == nil {
		t.Fatal("expected error for read of unwritten value")
	}
}

func TestValidateOK(t *testing.T) {
	h := history.History{
		{Process: 1, Type: history.Ok, Mops: []history.Mop{{Type: history.Append, Key: "x", App: 1}}},
		{Process: 2, Type: history.Ok, Mops: []history.Mop{{Type: history.Read, Key: "x", Read: []int{1}}}},
	}
	if err := Validate(h); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
