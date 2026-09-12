package history

type MopType int

const (
	Append MopType = 0
	Read   MopType = 1
)

type Mop struct {
	Type MopType
	Key  string
	App  int
	Read []int
}
