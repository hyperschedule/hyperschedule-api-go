package update

import (
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/db"
	"github.com/jackc/pgx/v4"
)

type SummaryCopy struct {
	Term            int64
	Staff           int64
	Course          int64
	Section         int64
	SectionStaff    int64
	SectionSchedule int64
}

func summaryCopyFromSeq(counts [6]int64) SummaryCopy {
	return SummaryCopy{
		Term:            counts[0],
		Staff:           counts[1],
		Course:          counts[2],
		Section:         counts[3],
		SectionStaff:    counts[4],
		SectionSchedule: counts[5],
	}
}

func (tbl table) copyFrom(tx *db.Tx, rows [][]interface{}) (int64, error) {
	n, err := tx.CopyFrom(
		pgx.Identifier{tbl.name},
		colNames(tbl.columns),
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return 0, err
	}
	if int64(len(rows)) != n {
		return n, fmt.Errorf("mismatch in # of copied rows: expect %d, actual %d", len(rows), n)
	}

	return n, nil
}

func (tbls *tables) copyFrom(tx *db.Tx, rs *rows) (SummaryCopy, error) {
	counts := [6]int64{}

	rowsOrder := rs.seq()
	for i, tbl := range tbls.seq() {
		count, err := tbl.copyFrom(tx, rowsOrder[i])
		if err != nil {
			return SummaryCopy{}, fmt.Errorf("failed to copy into %#v: %w", tbl.name, err)
		}
		counts[i] = count
	}

	return summaryCopyFromSeq(counts), nil
}
