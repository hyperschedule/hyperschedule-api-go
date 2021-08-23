package db

import (
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
	"github.com/jackc/pgx/v4"
)

type updateInfoTerm struct {
	Upserted int64
}

func (h *handle) updateTerms(terms map[string]*data.Term) (*updateInfoTerm, error) {
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
    , "start_date" date NOT NULL
    , "end_date"   date NOT NULL
    )
    ON COMMIT DROP;
  `); err != nil {
		return nil, fmt.Errorf("failed to create `term_tmp` table: %w", err)
	}

	if _, err := h.tx.CopyFrom(h.ctx, pgx.Identifier{"term_tmp"},
		[]string{"code", "semester", "start_date", "end_date"},
		pgx.CopyFromRows(rows),
	); err != nil {
		return nil, fmt.Errorf("failed to populate `term_tmp` table: %w", err)
	}

	tag, err := h.tx.Exec(h.ctx, `
    INSERT INTO "term" ( "code", "semester", "start_date", "end_date" )
    SELECT "code", "semester", "start_date", "end_date"
    FROM "term_tmp"
    ON CONFLICT ("code") DO UPDATE SET
        "semester" = EXCLUDED."semester"
      , "start_date" = EXCLUDED."start_date"
      , "end_date" = EXCLUDED."end_date"
    WHERE
      ( "term"."semester"
      , "term"."start_date"
      , "term"."end_date"
      )
      <>
      ( EXCLUDED."semester"
      , EXCLUDED."start_date"
      , EXCLUDED."end_date"
      );
  `)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert `term` entries: %w", err)
	}

	return &updateInfoTerm{
		Upserted: tag.RowsAffected(),
	}, nil
}
