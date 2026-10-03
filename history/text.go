package history

import "fmt"

var mopNames = map[MopType]string{Append: "append", Read: "r"}

var opNames = map[OpType]string{Invoke: "invoke", Ok: "ok", Fail: "fail", Info: "info"}

func (t MopType) MarshalText() ([]byte, error) {
	return []byte(mopNames[t]), nil
}

func (t *MopType) UnmarshalText(b []byte) error {
	return parseName(mopNames, string(b), t)
}

func (t OpType) MarshalText() ([]byte, error) {
	return []byte(opNames[t]), nil
}

func (t *OpType) UnmarshalText(b []byte) error {
	return parseName(opNames, string(b), t)
}

func parseName[T comparable](names map[T]string, s string, out *T) error {
	for k, v := range names {
		if v == s {
			*out = k
			return nil
		}
	}
	return fmt.Errorf("unknown type %q", s)
}
