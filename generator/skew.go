package generator

import (
	"math/rand"

	"elle-go/history"
)

func WriteSkew(cfg Config, seed int64) history.History {
	r := rand.New(rand.NewSource(seed))
	i := r.Intn(len(cfg.Keys))
	j := (i + 1 + r.Intn(len(cfg.Keys)-1)) % len(cfg.Keys)
	kx, ky := cfg.Keys[i], cfg.Keys[j]
	return history.History{
		{Process: 1, Type: history.Ok, Mops: []history.Mop{
			{Type: history.Read, Key: ky, Read: []int{}},
			{Type: history.Append, Key: kx, App: 1},
		}},
		{Process: 2, Type: history.Ok, Mops: []history.Mop{
			{Type: history.Read, Key: kx, Read: []int{}},
			{Type: history.Append, Key: ky, App: 2},
		}},
	}
}
