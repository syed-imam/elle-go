package generator

import (
	"math/rand"

	"github.com/syed-imam/elle-go/history"
)

func Invocations(cfg Config, seed int64) []history.Op {
	r := rand.New(rand.NewSource(seed))
	next := 1
	ops := make([]history.Op, 0, cfg.Txns)
	for i := 0; i < cfg.Txns; i++ {
		n := 1 + r.Intn(cfg.MaxMops)
		mops := make([]history.Mop, 0, n)
		for j := 0; j < n; j++ {
			key := cfg.Keys[r.Intn(len(cfg.Keys))]
			if r.Intn(2) == 0 {
				mops = append(mops, history.Mop{Type: history.Append, Key: key, App: next})
				next++
			} else {
				mops = append(mops, history.Mop{Type: history.Read, Key: key})
			}
		}
		ops = append(ops, history.Op{Type: history.Invoke, Mops: mops})
	}
	return ops
}
