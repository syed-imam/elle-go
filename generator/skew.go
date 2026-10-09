package generator

import (
	"math/rand"

	"elle-go/history"
)

func keyPair(cfg Config, seed int64) (string, string) {
	r := rand.New(rand.NewSource(seed))
	i := r.Intn(len(cfg.Keys))
	j := (i + 1 + r.Intn(len(cfg.Keys)-1)) % len(cfg.Keys)
	return cfg.Keys[i], cfg.Keys[j]
}

func ReadSkew(cfg Config, seed int64) history.History {
	kx, ky := keyPair(cfg, seed)
	return history.History{
		{Process: 1, Type: history.Ok, Mops: []history.Mop{
			{Type: history.Read, Key: kx, Read: []int{}},
			{Type: history.Read, Key: ky, Read: []int{2}},
		}},
		{Process: 2, Type: history.Ok, Mops: []history.Mop{
			{Type: history.Append, Key: kx, App: 1},
			{Type: history.Append, Key: ky, App: 2},
		}},
	}
}

func WriteSkew(cfg Config, seed int64) history.History {
	kx, ky := keyPair(cfg, seed)
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

func CircularInformationFlow(cfg Config, seed int64) history.History {
	kx, ky := keyPair(cfg, seed)
	return history.History{
		{Process: 1, Type: history.Ok, Mops: []history.Mop{
			{Type: history.Append, Key: kx, App: 1},
			{Type: history.Read, Key: ky, Read: []int{2}},
		}},
		{Process: 2, Type: history.Ok, Mops: []history.Mop{
			{Type: history.Append, Key: ky, App: 2},
			{Type: history.Read, Key: kx, Read: []int{1}},
		}},
	}
}

func WriteCycle(cfg Config, seed int64) history.History {
	kx, ky := keyPair(cfg, seed)
	return history.History{
		{Process: 1, Type: history.Ok, Mops: []history.Mop{
			{Type: history.Append, Key: kx, App: 1},
			{Type: history.Append, Key: ky, App: 1},
		}},
		{Process: 2, Type: history.Ok, Mops: []history.Mop{
			{Type: history.Append, Key: kx, App: 2},
			{Type: history.Append, Key: ky, App: 2},
		}},
		{Process: 3, Type: history.Ok, Mops: []history.Mop{
			{Type: history.Read, Key: kx, Read: []int{1, 2}},
			{Type: history.Read, Key: ky, Read: []int{2, 1}},
		}},
	}
}

func DivergentReads(cfg Config, seed int64) history.History {
	kx, _ := keyPair(cfg, seed)
	return history.History{
		{Process: 1, Type: history.Ok, Mops: []history.Mop{{Type: history.Append, Key: kx, App: 1}}},
		{Process: 2, Type: history.Ok, Mops: []history.Mop{{Type: history.Append, Key: kx, App: 2}}},
		{Process: 3, Type: history.Ok, Mops: []history.Mop{{Type: history.Read, Key: kx, Read: []int{1, 2}}}},
		{Process: 4, Type: history.Ok, Mops: []history.Mop{{Type: history.Read, Key: kx, Read: []int{2, 1}}}},
	}
}

func LostOwnAppend(cfg Config, seed int64) history.History {
	kx, _ := keyPair(cfg, seed)
	return history.History{
		{Process: 1, Type: history.Ok, Mops: []history.Mop{{Type: history.Append, Key: kx, App: 1}}},
		{Process: 2, Type: history.Ok, Mops: []history.Mop{
			{Type: history.Append, Key: kx, App: 2},
			{Type: history.Read, Key: kx, Read: []int{1}},
		}},
	}
}
