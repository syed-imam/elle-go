package generator

import (
	"testing"

	"github.com/syed-imam/elle-go/anomaly"
	"github.com/syed-imam/elle-go/differential"
)

func TestSerialIsSerializable(t *testing.T) {
	cfg := Config{Keys: []string{"x", "y", "z"}, Txns: 20, MaxMops: 4}
	for seed := int64(0); seed < 200; seed++ {
		h := Serial(cfg, seed)
		if v, err := differential.Verdict(h); err != nil || v.Anomaly != anomaly.None {
			t.Fatalf("seed %d: expected None, got %v (cycle %v, err %v)", seed, v.Anomaly, v.Cycle, err)
		}
	}
}
