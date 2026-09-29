package generator

import (
	"testing"

	"elle-go/anomaly"
	"elle-go/differential"
)

func TestWriteSkewIsG2(t *testing.T) {
	cfg := Config{Keys: []string{"x", "y", "z"}}
	for seed := int64(0); seed < 200; seed++ {
		h := WriteSkew(cfg, seed)
		if v := differential.Verdict(h); v.Anomaly != anomaly.G2 {
			t.Fatalf("seed %d: expected G2, got %v (cycle %v)", seed, v.Anomaly, v.Cycle)
		}
	}
}
