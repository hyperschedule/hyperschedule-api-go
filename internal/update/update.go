package update

import (
	"context"
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
	"github.com/MuddCreates/hyperschedule-api-go/internal/db"
	"log"
)

type Summary struct {
	Copy  SummaryCopy
	Apply SummaryApply
}

func Update(ctx context.Context, pool *db.Connection, d *data.Data) (Summary, error) {
	log.Printf("initiating transaction")
	tx, err := pool.Begin(ctx)
	if err != nil {
		return Summary{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Cleanup()

	log.Printf("creating temporary tables")
	if err := tmpTables.create(tx); err != nil {
		return Summary{}, fmt.Errorf("failed to create temporary tables: %w", err)
	}

	log.Printf("populating temporary tables")
	summaryCopy, err := tmpTables.copyFrom(tx, rowsFrom(d))
	if err != nil {
		return Summary{}, fmt.Errorf("failed to populate temporary tables: %w", err)
	}

	log.Printf("applying modifications to persistent tables")
	summaryApply, err := tmpTables.apply(tx)
	if err != nil {
		return Summary{}, fmt.Errorf("failed to apply temporary table changes: %w", err)
	}

	log.Printf("committing changes")
	if err := tx.Commit(); err != nil {
		return Summary{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("done")
	return Summary{
		Copy:  summaryCopy,
		Apply: summaryApply,
	}, nil
}
