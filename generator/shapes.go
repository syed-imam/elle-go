package generator

import (
	"math/rand"

	"github.com/syed-imam/elle-go/history"
)

func FromShapes(shapes []history.Op, txns int, seed int64) []history.Op {
	r := rand.New(rand.NewSource(seed))
	next := 1
	ops := make([]history.Op, 0, txns)
	for i := 0; i < txns; i++ {
		shape := shapes[r.Intn(len(shapes))]
		mops := make([]history.Mop, len(shape.Mops))
		for j, m := range shape.Mops {
			mops[j] = history.Mop{Type: m.Type, Key: m.Key}
			if m.Type == history.Append {
				mops[j].App = next
				next++
			}
		}
		ops = append(ops, history.Op{Type: history.Invoke, Mops: mops})
	}
	return ops
}
