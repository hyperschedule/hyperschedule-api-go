package update

import (
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/db"
	"github.com/jackc/pgx/v4"
	"strings"
)

func (ts tables) create(tx *db.Tx) error {
	results := tx.Batch(ts.batchCreate())
	defer results.Close()

	for _, tbl := range ts.seq() {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf(
				"failed to create temporary table %#v: %w",
				tbl.name, err,
			)
		}
	}
	if err := results.Close(); err != nil {
		return fmt.Errorf("failed to close temporary table batch: %w", err)
	}
	return nil
}

func (tbl table) sqlCreate() string {
	sqlSpecs := make([]string, 0, len(tbl.columns)+len(tbl.constraints))

	for _, spec := range tbl.columns {
		sqlSpecs = append(sqlSpecs, spec.sql())
	}
	for _, spec := range tbl.constraints {
		sqlSpecs = append(sqlSpecs, string(spec))
	}

	return fmt.Sprintf(
		"CREATE TEMPORARY TABLE %s (%s) ON COMMIT DROP;",
		pgx.Identifier{tbl.name}.Sanitize(),
		strings.Join(sqlSpecs, ","),
	)
}

func (spec columnSpec) sql() string {
	return fmt.Sprintf(
		"%s %s", pgx.Identifier{spec.name}.Sanitize(), spec.props,
	)
}

func pkey(cols ...string) constraintSpec {
	colsSanitized := make([]string, len(cols))
	for i, col := range cols {
		colsSanitized[i] = pgx.Identifier{col}.Sanitize()
	}
	return constraintSpec(fmt.Sprintf(
		"PRIMARY KEY (%s)", strings.Join(colsSanitized, ","),
	))
}

func (ts tables) batchCreate() *pgx.Batch {
	batch := &pgx.Batch{}
	for _, tbl := range ts.seq() {
		batch.Queue(tbl.sqlCreate())
	}
	return batch
}
