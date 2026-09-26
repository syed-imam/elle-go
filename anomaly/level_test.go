package anomaly

import "testing"

func TestRequires(t *testing.T) {
	cases := []struct {
		a    Anomaly
		want Level
	}{
		{None, ReadUncommitted},
		{G0, ReadUncommitted},
		{G1c, ReadCommitted},
		{GSingle, SnapshotIsolation},
		{G2, Serializable},
	}
	for _, c := range cases {
		if got := Requires(c.a); got != c.want {
			t.Errorf("Requires(%v) = %v, want %v", c.a, got, c.want)
		}
	}
}
