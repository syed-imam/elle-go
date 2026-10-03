package history

import (
	"encoding/json"
	"testing"
)

func TestOpJSONRoundTrip(t *testing.T) {
	in := `{"process":1,"type":"ok","mops":[{"type":"append","key":"x","app":1,"read":null},{"type":"r","key":"x","app":0,"read":[1]}]}`
	var op Op
	if err := json.Unmarshal([]byte(in), &op); err != nil {
		t.Fatal(err)
	}
	if op.Type != Ok || op.Mops[0].Type != Append || op.Mops[1].Type != Read {
		t.Fatalf("decoded wrong types: %+v", op)
	}
	out, err := json.Marshal(op)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != in {
		t.Fatalf("round trip mismatch:\n got %s\nwant %s", out, in)
	}
}

func TestUnknownTypeRejected(t *testing.T) {
	var m Mop
	if err := json.Unmarshal([]byte(`{"type":"write"}`), &m); err == nil {
		t.Fatal("expected error for unknown mop type")
	}
}
