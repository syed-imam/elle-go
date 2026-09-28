package differential

import (
	"os"
	"testing"

	"elle-go/anomaly"
	"elle-go/fixtures"
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
