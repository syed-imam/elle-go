package differential

import (
	"github.com/syed-imam/elle-go/anomaly"
	"github.com/syed-imam/elle-go/checker"
	"github.com/syed-imam/elle-go/history"
)

func Verdict(h history.History) (anomaly.Verdict, error) {
	return checker.Check(h)
}
