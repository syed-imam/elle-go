package differential

import (
	"testing"

	"elle-go/anomaly"
)

func TestTypeAgrees(t *testing.T) {
	if !TypeAgrees(anomaly.G2, []string{"G2-item"}) {
		t.Error("G2 should agree with [G2-item]")
	}
	if !TypeAgrees(anomaly.None, nil) {
		t.Error("None should agree with no elle types")
	}
	if TypeAgrees(anomaly.G2, []string{"G1c"}) {
		t.Error("G2 should not agree with [G1c]")
	}
	if TypeAgrees(anomaly.None, []string{"G2-item"}) {
		t.Error("None should not agree when elle reports a type")
	}
}

func TestTypeAgreesIncompatibleOrder(t *testing.T) {
	if !TypeAgrees(anomaly.IncompatibleOrder, []string{"G0", "incompatible-order"}) {
		t.Error("IncompatibleOrder should agree with [G0 incompatible-order]")
	}
}

func TestTypeAgreesInternal(t *testing.T) {
	if !TypeAgrees(anomaly.Internal, []string{"G-single-item", "internal"}) {
		t.Error("Internal should agree with [G-single-item internal]")
	}
}

func TestTypeAgreesLostUpdate(t *testing.T) {
	if !TypeAgrees(anomaly.LostUpdate, []string{"lost-update"}) {
		t.Error("LostUpdate should agree with [lost-update]")
	}
}
