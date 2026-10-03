package history

type MopType int

const (
	Append MopType = 0
	Read   MopType = 1
)

type Mop struct {
	Type MopType `json:"type"`
	Key  string  `json:"key"`
	App  int     `json:"app"`
	Read []int   `json:"read"`
}
