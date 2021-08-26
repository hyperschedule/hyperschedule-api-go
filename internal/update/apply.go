package update

import (
	"fmt"
	_ "github.com/MuddCreates/hyperschedule-api-go/internal/data"
	"github.com/MuddCreates/hyperschedule-api-go/internal/db"
	"github.com/jackc/pgx/v4"
)

type SummaryApplyStats struct {
	Upsert int64
	Delete int64
}

type SummaryApply struct {
	Term            SummaryApplyStats
	Staff           SummaryApplyStats
	Course          SummaryApplyStats
	Section         SummaryApplyStats
	SectionStaff    SummaryApplyStats
	SectionSchedule SummaryApplyStats
}

// TODO add generics to go :))))))))))))))))) :(
func summaryApplyFromSeq(stats [6]SummaryApplyStats) SummaryApply {
	return SummaryApply{
		Term:            stats[0],
		Staff:           stats[1],
		Course:          stats[2],
		Section:         stats[3],
		SectionStaff:    stats[4],
		SectionSchedule: stats[5],
	}
}

func (ts tables) apply(tx *db.Tx) (SummaryApply, error) {
	results := tx.Batch(ts.batchApply())
	defer results.Close()

	next := func(s sqlApply) (SummaryApplyStats, error) {
		tagUpsert, err := results.Exec()
		if err != nil {
			return SummaryApplyStats{}, fmt.Errorf("upsert: %w", err)
		}

		tagDelete, err := results.Exec()
		if err != nil {
			return SummaryApplyStats{}, fmt.Errorf("delete: %w", err)
		}

		return SummaryApplyStats{
			Upsert: tagUpsert.RowsAffected(),
			Delete: tagDelete.RowsAffected(),
		}, nil
	}

	stats := [6]SummaryApplyStats{}
	for i, tbl := range ts.seq() {
		s, err := next(tbl.sqlApply)
		if err != nil {
			return SummaryApply{}, fmt.Errorf("table %q: %w", tbl.name, err)
		}
		stats[i] = s
	}

	return summaryApplyFromSeq(stats), nil
}

func (ts tables) batchApply() *pgx.Batch {
	batch := &pgx.Batch{}
	for _, tbl := range ts.seq() {
		batch.Queue(tbl.sqlApply.upsert)
		batch.Queue(tbl.sqlApply.delete)
	}
	return batch
}
