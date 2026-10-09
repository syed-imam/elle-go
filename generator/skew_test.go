package generator

import (
	"testing"

	"elle-go/anomaly"
	"elle-go/differential"
)

func TestReadSkewIsGSingle(t *testing.T) {
	cfg := Config{Keys: []string{"x", "y", "z"}}
	for seed := int64(0); seed < 200; seed++ {
		h := ReadSkew(cfg, seed)
		if v, err := differential.Verdict(h); err != nil || v.Anomaly != anomaly.GSingle {
			t.Fatalf("seed %d: expected G-single, got %v (cycle %v, err %v)", seed, v.Anomaly, v.Cycle, err)
		}
	}
}

func TestWriteSkewIsG2(t *testing.T) {
	cfg := Config{Keys: []string{"x", "y", "z"}}
	for seed := int64(0); seed < 200; seed++ {
		h := WriteSkew(cfg, seed)
		if v, err := differential.Verdict(h); err != nil || v.Anomaly != anomaly.G2 {
			t.Fatalf("seed %d: expected G2, got %v (cycle %v, err %v)", seed, v.Anomaly, v.Cycle, err)
		}
	}
}

func TestCircularInformationFlowIsG1c(t *testing.T) {
	cfg := Config{Keys: []string{"x", "y", "z"}}
	for seed := int64(0); seed < 200; seed++ {
		h := CircularInformationFlow(cfg, seed)
		if v, err := differential.Verdict(h); err != nil || v.Anomaly != anomaly.G1c {
			t.Fatalf("seed %d: expected G1c, got %v (cycle %v, err %v)", seed, v.Anomaly, v.Cycle, err)
		}
	}
}

func TestWriteCycleIsG0(t *testing.T) {
	cfg := Config{Keys: []string{"x", "y", "z"}}
	for seed := int64(0); seed < 200; seed++ {
		h := WriteCycle(cfg, seed)
		if v, err := differential.Verdict(h); err != nil || v.Anomaly != anomaly.G0 {
			t.Fatalf("seed %d: expected G0, got %v (cycle %v, err %v)", seed, v.Anomaly, v.Cycle, err)
		}
	}
}

func TestDivergentReadsIsIncompatibleOrder(t *testing.T) {
	cfg := Config{Keys: []string{"x", "y", "z"}}
	for seed := int64(0); seed < 200; seed++ {
		h := DivergentReads(cfg, seed)
		if v, err := differential.Verdict(h); err != nil || v.Anomaly != anomaly.IncompatibleOrder {
			t.Fatalf("seed %d: expected incompatible-order, got %v (err %v)", seed, v.Anomaly, err)
		}
	}
}

func TestLostOwnAppendIsInternal(t *testing.T) {
	cfg := Config{Keys: []string{"x", "y", "z"}}
	for seed := int64(0); seed < 200; seed++ {
		h := LostOwnAppend(cfg, seed)
		if v, err := differential.Verdict(h); err != nil || v.Anomaly != anomaly.Internal {
			t.Fatalf("seed %d: expected internal, got %v (err %v)", seed, v.Anomaly, err)
		}
	}
}
