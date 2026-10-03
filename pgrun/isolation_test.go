package pgrun

import (
	"database/sql"
	"testing"
)

func TestParseIsolation(t *testing.T) {
	if iso, err := ParseIsolation("repeatable-read"); err != nil || iso != sql.LevelRepeatableRead {
		t.Fatalf("got %v, %v", iso, err)
	}
	if _, err := ParseIsolation("snapshot"); err == nil {
		t.Fatal("expected error for unknown isolation")
	}
}
