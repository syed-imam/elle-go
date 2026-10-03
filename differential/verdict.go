package differential

import (
	"elle-go/anomaly"
	"elle-go/checker"
	"elle-go/history"
)

func Verdict(h history.History) (anomaly.Verdict, error) {
	return checker.Check(h)
}
