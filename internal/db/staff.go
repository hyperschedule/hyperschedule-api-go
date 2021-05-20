package db

import (
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
)

type updateInfoStaff struct {
	Count int64
}

func updateStaff(h *handle, staff map[string]data.Name) (*updateInfoStaff, error) {
	rows := make([][]interface{}, 0)
	for lingkId, name := range staff {
		rows = append(rows, []interface{}{lingkId, name.First, name.Last})
	}
	result, err := batchInsert(h, &batchOpts{
		tmpName: "staff_tmp",
		tmpCols: []string{"lingk_id", "first_name", "last_name"},
		sqlTmpDefs: `
      "lingk_id" text PRIMARY KEY
    , "first_name" text NOT NULL
    , "last_name" text NOT NULL
    `,
		sqlInsertSelect: `
      INSERT INTO "staff" ("lingk_id", "first_name", "last_name")
      SELECT "lingk_id", "first_name", "last_name" FROM "staff_tmp"
      ON CONFLICT ("lingk_id") DO UPDATE SET
        "first_name" = EXCLUDED."first_name"
      , "last_name" = EXCLUDED."last_name"
      WHERE 
        "staff"."first_name" <> EXCLUDED."first_name"
        OR "staff"."last_name" <> EXCLUDED."last_name"
    `,
		rows: rows,
	})
	if err != nil {
		return nil, err
	}
	return &updateInfoStaff{
		Count: result.RowsAffected(),
	}, nil
}
