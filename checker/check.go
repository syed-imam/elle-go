package checker

import (
	"github.com/syed-imam/elle-go/anomaly"
	"github.com/syed-imam/elle-go/graph"
	"github.com/syed-imam/elle-go/history"
	"github.com/syed-imam/elle-go/listappend"
)

func Check(h history.History) (anomaly.Verdict, error) {
	if err := listappend.Validate(h); err != nil {
		return anomaly.Verdict{}, err
	}
	if c, ok := listappend.IncompatibleOrder(h); ok {
		return anomaly.IncompatibleOrderVerdict(c.Key, c.Shorter, c.Longer), nil
	}
	g := graph.New()
	for _, e := range listappend.Dependencies(h) {
		g.AddEdge(e.From, e.To, toRel(e.Type))
	}
	v := anomaly.Check(g)
	if v.Anomaly != anomaly.None && v.Anomaly != anomaly.G2 {
		return v, nil
	}
	if c, ok := listappend.Internal(h); ok {
		return anomaly.InternalVerdict(c.Op, c.Key, c.Expected, c.Prefixed, c.Read), nil
	}
	if c, ok := listappend.LostUpdate(h); ok {
		return anomaly.LostUpdateVerdict(c.Key, c.Value, c.Ops), nil
	}
	return v, nil
}

func toRel(t listappend.EdgeType) graph.Rel {
	switch t {
	case listappend.WR:
		return graph.WR
	case listappend.RW:
		return graph.RW
	default:
		return graph.WW
	}
}
