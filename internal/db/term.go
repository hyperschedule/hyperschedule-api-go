package db

import (
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
	"github.com/jackc/pgx/v4"
)

func (h *handle) updateTerms(terms map[string]*data.Term) (UpdateSummary, error) {
	rows := [][]interface{}{}
	for code, term := range terms {
		rows = append(rows, []interface{}{
			code, term.Semester, term.Start.ToTime(), term.End.ToTime(),
		})
	}

	if _, err := h.tx.Exec(h.ctx, `
    CREATE TEMPORARY TABLE "term_tmp"
    ( "code"       text PRIMARY KEY
    , "semester"   text NOT NULL
    , "date_start" date NOT NULL
    , "date_end"   date NOT NULL
    )
    ON COMMIT DROP;
  `); err != nil {
		return UpdateSummary{}, fmt.Errorf("failed to create `term_tmp` table: %w", err)
	}

	countCopy, err := h.tx.CopyFrom(h.ctx, pgx.Identifier{"term_tmp"},
		[]string{"code", "semester", "date_start", "date_end"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return UpdateSummary{}, fmt.Errorf("failed to populate `term_tmp` table: %w", err)
	}

	tag, err := h.tx.Exec(h.ctx, `
    INSERT INTO "term" ( "code", "semester", "date_start", "date_end" )
    SELECT "code", "semester", "date_start", "date_end"
    FROM "term_tmp"
    ON CONFLICT ("code") DO UPDATE SET
        "semester" = EXCLUDED."semester"
      , "date_start" = EXCLUDED."date_start"
      , "date_end" = EXCLUDED."date_end"
    WHERE
      ( "term"."semester"
      , "term"."date_start"
      , "term"."date_end"
      )
      <>
      ( EXCLUDED."semester"
      , EXCLUDED."date_start"
      , EXCLUDED."date_end"
      );
  `)
	if err != nil {
		return UpdateSummary{}, fmt.Errorf("failed to upsert `term` entries: %w", err)
	}

	return UpdateSummary{
		Copied:   countCopy,
		Upserted: tag.RowsAffected(),
	}, nil
}
