package bbolt

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"time"

	bolt "go.etcd.io/bbolt"
)

// Connect to a bbolt database
func Connect(filepath string) (*bolt.DB, error) {
	dirName := path.Dir(filepath)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			return nil, fmt.Errorf("creating dir: %w", merr)
		}
	}

	db, err := bolt.Open(filepath, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	return db, nil
}

func Stats(ctx context.Context, db *bolt.DB, t time.Duration, logger *slog.Logger) {
	go func() {
		// Grab the initial stats.
		prev := db.Stats()
		ticker := time.NewTicker(t)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Grab the current stats and diff them.
				stats := db.Stats()
				diff := stats.Sub(&prev)

				logger.InfoContext(ctx,
					"bolt db stats",
					"free_page_n", diff.FreePageN,
					"pending_page_n", diff.PendingPageN,
					"free_alloc", diff.FreeAlloc,
					"freelist_inuse", diff.FreelistInuse,
					"tx_n", diff.TxN,
					"open_tx_n", diff.OpenTxN,
				)

				// Save stats for the next loop.
				prev = stats
			}
		}
	}()
}
