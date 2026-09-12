package history

type OpType int

const (
	Invoke OpType = 0
	Ok     OpType = 1
	Fail   OpType = 2
	Info   OpType = 3
)

type Op struct {
	Process int
	Type    OpType
	Mops    []Mop
}
