package pgrun

import (
	"database/sql"
	"fmt"
)

var isolations = map[string]sql.IsolationLevel{
	"read-committed":  sql.LevelReadCommitted,
	"repeatable-read": sql.LevelRepeatableRead,
	"serializable":    sql.LevelSerializable,
}

func ParseIsolation(s string) (sql.IsolationLevel, error) {
	if iso, ok := isolations[s]; ok {
		return iso, nil
	}
	return 0, fmt.Errorf("unknown isolation %q (want read-committed, repeatable-read, or serializable)", s)
}
