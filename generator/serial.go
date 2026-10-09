package generator

import (
	"math/rand"

	"github.com/syed-imam/elle-go/history"
)

type Config struct {
	Keys    []string
	Txns    int
	MaxMops int
}

func Serial(cfg Config, seed int64) history.History {
	r := rand.New(rand.NewSource(seed))
	state := map[string][]int{}
	next := 1
	h := make(history.History, 0, cfg.Txns)
	for i := 0; i < cfg.Txns; i++ {
		n := 1 + r.Intn(cfg.MaxMops)
		mops := make([]history.Mop, 0, n)
		for j := 0; j < n; j++ {
			key := cfg.Keys[r.Intn(len(cfg.Keys))]
			if r.Intn(2) == 0 {
				state[key] = append(state[key], next)
				mops = append(mops, history.Mop{Type: history.Append, Key: key, App: next})
				next++
			} else {
				mops = append(mops, history.Mop{Type: history.Read, Key: key, Read: append([]int{}, state[key]...)})
			}
		}
		h = append(h, history.Op{Process: i + 1, Type: history.Ok, Mops: mops})
	}
	return h
}
