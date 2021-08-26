package update

import (
	"context"
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
	"github.com/MuddCreates/hyperschedule-api-go/internal/db"
)

type Summary struct {
	Copy  SummaryCopy
	Apply SummaryApply
}

func Update(ctx context.Context, pool *db.Connection, d *data.Data) (Summary, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return Summary{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Cleanup()

	if err := tmpTables.create(tx); err != nil {
		return Summary{}, fmt.Errorf("failed to create temporary tables: %w", err)
	}

	summaryCopy, err := tmpTables.copyFrom(tx, rowsFrom(d))
	if err != nil {
		return Summary{}, fmt.Errorf("failed to populate temporary tables: %w", err)
	}

	summaryApply, err := tmpTables.apply(tx)
	if err != nil {
		return Summary{}, fmt.Errorf("failed to apply temporary table changes: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return Summary{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return Summary{
		Copy:  summaryCopy,
		Apply: summaryApply,
	}, nil
}
