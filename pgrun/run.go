package pgrun

import (
	"context"
	"database/sql"
	"sync"

	"elle-go/history"
)

func Run(ctx context.Context, db *sql.DB, iso sql.IsolationLevel, ops []history.Op, workers int) history.History {
	jobs := make(chan history.Op)
	var mu sync.Mutex
	var wg sync.WaitGroup
	h := history.History{}
	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go func(process int) {
			defer wg.Done()
			for op := range jobs {
				op.Process = process
				done := Exec(ctx, db, iso, op)
				mu.Lock()
				h = append(h, done)
				mu.Unlock()
			}
		}(w)
	}
	for _, op := range ops {
		jobs <- op
	}
	close(jobs)
	wg.Wait()
	return h
}
