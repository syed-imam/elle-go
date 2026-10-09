package differential

import "elle-go/anomaly"

func elleName(a anomaly.Anomaly) string {
	switch a {
	case anomaly.G0:
		return "G0"
	case anomaly.G1c:
		return "G1c"
	case anomaly.GSingle:
		return "G-single-item"
	case anomaly.G2:
		return "G2-item"
	case anomaly.IncompatibleOrder:
		return "incompatible-order"
	case anomaly.Internal:
		return "internal"
	default:
		return ""
	}
}

func TypeAgrees(ours anomaly.Anomaly, elle []string) bool {
	want := elleName(ours)
	if want == "" {
		return len(elle) == 0
	}
	for _, t := range elle {
		if t == want {
			return true
		}
	}
	return false
}
