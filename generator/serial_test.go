package generator

import (
	"testing"

	"elle-go/anomaly"
	"elle-go/differential"
)

func TestSerialIsSerializable(t *testing.T) {
	cfg := Config{Keys: []string{"x", "y", "z"}, Txns: 20, MaxMops: 4}
	for seed := int64(0); seed < 200; seed++ {
		h := Serial(cfg, seed)
		if v := differential.Verdict(h); v.Anomaly != anomaly.None {
			t.Fatalf("seed %d: expected None, got %v (cycle %v)", seed, v.Anomaly, v.Cycle)
		}
	}
}
