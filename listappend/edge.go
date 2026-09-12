package listappend

type EdgeType int

const (
	WW EdgeType = 0
	WR EdgeType = 1
	RW EdgeType = 2
)

type Edge struct {
	From int
	To   int
	Type EdgeType
}
