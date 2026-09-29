package differential

import (
	"os"
	"testing"

	"elle-go/anomaly"
	"elle-go/fixtures"
	"elle-go/generator"
	"elle-go/history"
)

func TestParityFixtures(t *testing.T) {
	if os.Getenv("ELLE_PARITY") == "" {
		t.Skip("set ELLE_PARITY=1 to run the Elle differential test (slow: JVM per history)")
	}
	if !ElleAvailable() {
		t.Skip("lein or ../../elle-upstream not available; skipping Elle differential test")
	}
	cases := []struct {
		name string
		h    history.History
	}{
		{"writeskew", fixtures.WriteSkew()},
		{"serializable", fixtures.Serializable()},
	}
	for _, c := range cases {
		elle, err := RunElle(c.h)
		if err != nil {
			t.Fatalf("%s: RunElle: %v", c.name, err)
		}
		ours := Verdict(c.h)
		ourViolation := ours.Anomaly != anomaly.None
		if ourViolation != elle.Violation() {
			t.Errorf("%s: violation disagreement: elle.valid?=%q (viol=%v %v) vs ours=%v",
				c.name, elle.ValidField, elle.Violation(), elle.AnomalyTypes, ours)
			continue
		}
		t.Logf("%s: agree — elle viol=%v %v, ours=%v", c.name, elle.Violation(), elle.AnomalyTypes, ours)
	}
}

func TestParityGenerated(t *testing.T) {
	if os.Getenv("ELLE_PARITY") == "" {
		t.Skip("set ELLE_PARITY=1 to run the Elle differential test (slow: JVM per history)")
	}
	if !ElleAvailable() {
		t.Skip("lein or ../../elle-upstream not available; skipping Elle differential test")
	}
	cfg := generator.Config{Keys: []string{"x", "y", "z"}, Txns: 10, MaxMops: 3}
	for seed := int64(0); seed < 5; seed++ {
		h := generator.Serial(cfg, seed)
		elle, err := RunElle(h)
		if err != nil {
			t.Fatalf("seed %d: RunElle: %v", seed, err)
		}
		ours := Verdict(h)
		ourViolation := ours.Anomaly != anomaly.None
		if ourViolation != elle.Violation() {
			t.Errorf("seed %d: disagree: elle valid?=%q %v vs ours=%v",
				seed, elle.ValidField, elle.AnomalyTypes, ours)
			continue
		}
		if ourViolation {
			t.Errorf("seed %d: serial history flagged as violation (elle=%v ours=%v)",
				seed, elle.AnomalyTypes, ours)
			continue
		}
		t.Logf("seed %d: agree — both valid", seed)
	}
}

func TestParityWriteSkew(t *testing.T) {
	if os.Getenv("ELLE_PARITY") == "" {
		t.Skip("set ELLE_PARITY=1 to run the Elle differential test (slow: JVM per history)")
	}
	if !ElleAvailable() {
		t.Skip("lein or ../../elle-upstream not available; skipping Elle differential test")
	}
	cfg := generator.Config{Keys: []string{"x", "y", "z"}}
	for seed := int64(0); seed < 5; seed++ {
		h := generator.WriteSkew(cfg, seed)
		elle, err := RunElle(h)
		if err != nil {
			t.Fatalf("seed %d: RunElle: %v", seed, err)
		}
		ours := Verdict(h)
		if ours.Anomaly == anomaly.None || !elle.Violation() {
			t.Errorf("seed %d: expected both to flag a violation: elle valid?=%q %v vs ours=%v",
				seed, elle.ValidField, elle.AnomalyTypes, ours)
			continue
		}
		if !TypeAgrees(ours.Anomaly, elle.AnomalyTypes) {
			t.Errorf("seed %d: type disagreement: ours=%v vs elle=%v",
				seed, ours.Anomaly, elle.AnomalyTypes)
			continue
		}
		t.Logf("seed %d: agree — ours=%v, elle=%v", seed, ours.Anomaly, elle.AnomalyTypes)
	}
}
