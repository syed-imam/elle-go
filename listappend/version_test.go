package listappend

import (
	"slices"
	"testing"

	"elle-go/fixtures"
)

func TestVersionOrder(t *testing.T) {
	orders := VersionOrder(fixtures.Serializable())

	if got, want := orders["x"], []int{1, 2}; !slices.Equal(got, want) {
		t.Errorf("x order = %v, want %v", got, want)
	}
	if got := orders["y"]; len(got) != 0 {
		t.Errorf("y order = %v, want empty", got)
	}
}
