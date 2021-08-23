package db

import (
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
	"github.com/jackc/pgx/v4"
)

func updateStaff(h *handle, staff map[string]data.Name) (UpdateSummary, error) {
	rows := make([][]interface{}, 0)
	for lingkId, name := range staff {
		rows = append(rows, []interface{}{lingkId, name.First, name.Last})
	}

	if _, err := h.tx.Exec(h.ctx, `
    CREATE TEMPORARY TABLE "staff_tmp"
      ( "lingk_id" text PRIMARY KEY
      , "first_name" text NOT NULL
      , "last_name" text NOT NULL
      )
    ON COMMIT DROP;
  `); err != nil {
		return UpdateSummary{}, fmt.Errorf("failed to create `staff_tmp`: %w", err)
	}

	countCopy, err := h.tx.CopyFrom(h.ctx, pgx.Identifier{"staff_tmp"},
		[]string{"lingk_id", "first_name", "last_name"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return UpdateSummary{}, fmt.Errorf("failed to populate `staff_tmp`: %w", err)
	}

	tagUpsert, err := h.tx.Exec(h.ctx, `
    INSERT INTO "staff" ("lingk_id", "first_name", "last_name")
    SELECT "lingk_id", "first_name", "last_name" FROM "staff_tmp"
    ON CONFLICT ("lingk_id") DO UPDATE SET
      "first_name" = EXCLUDED."first_name"
    , "last_name" = EXCLUDED."last_name"
    WHERE
      "staff"."first_name" <> EXCLUDED."first_name"
      OR "staff"."last_name" <> EXCLUDED."last_name"
  `)
	if err != nil {
		return UpdateSummary{}, fmt.Errorf("failed to upsert `staff` entries: %w", err)
	}

	return UpdateSummary{
		Copied:   countCopy,
		Upserted: tagUpsert.RowsAffected(),
	}, nil
}
