package differential

import (
	"fmt"
	"strings"

	"elle-go/history"
)

func ToEDN(h history.History) string {
	var b strings.Builder
	for _, op := range h {
		writeOp(&b, op, ":invoke", true)
		writeOp(&b, op, okType(op), false)
	}
	return b.String()
}

func okType(op history.Op) string {
	switch op.Type {
	case history.Fail:
		return ":fail"
	case history.Info:
		return ":info"
	default:
		return ":ok"
	}
}

func writeOp(b *strings.Builder, op history.Op, typ string, invoke bool) {
	fmt.Fprintf(b, "{:process %d :type %s :value [", op.Process, typ)
	for i, m := range op.Mops {
		if i > 0 {
			b.WriteByte(' ')
		}
		writeMop(b, m, invoke)
	}
	b.WriteString("]}\n")
}

func writeMop(b *strings.Builder, m history.Mop, invoke bool) {
	if m.Type == history.Append {
		fmt.Fprintf(b, "[:append %q %d]", m.Key, m.App)
		return
	}
	if invoke || len(m.Read) == 0 {
		fmt.Fprintf(b, "[:r %q nil]", m.Key)
		return
	}
	fmt.Fprintf(b, "[:r %q [", m.Key)
	for i, v := range m.Read {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(b, "%d", v)
	}
	b.WriteString("]]")
}
