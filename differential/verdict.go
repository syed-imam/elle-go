package differential

import (
	"elle-go/anomaly"
	"elle-go/checker"
	"elle-go/history"
)

func Verdict(h history.History) anomaly.Verdict {
	return checker.Check(h)
}
