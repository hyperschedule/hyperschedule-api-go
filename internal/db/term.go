package db

import (
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
)

type updateInfoTerm struct {
	count int64
}

func updateTerms(h *handle, terms map[string]*data.Term) (*updateInfoTerm, error) {
	rows := [][]interface{}{}
	for code, term := range terms {
		rows = append(rows, []interface{}{
			code, term.Semester, term.Start.ToTime(), term.End.ToTime(),
		})
	}

	result, err := batchInsert(h, &batchOpts{
		tmpName: "term_tmp",
		tmpCols: []string{"code", "semester", "start_date", "end_date"},
		sqlTmpDefs: `
      "code" text PRIMARY KEY
    , "semester" text NOT NULL
    , "start_date" date NOT NULL
    , "end_date" date NOT NULL
    `,
		sqlInsertSelect: `
      INSERT INTO "term" ("code", "semester", "start_date", "end_date")
      SELECT "code", "semester", "start_date", "end_date" FROM "term_tmp"
      ON CONFLICT ("code") DO UPDATE SET
        "semester" = EXCLUDED."semester"
      , "start_date" = EXCLUDED."start_date"
      , "end_date" = EXCLUDED."end_date"
      WHERE NOT (
        "term"."semester" = EXCLUDED."semester"
        AND "term"."start_date" = EXCLUDED."start_date"
        AND "term"."end_date" = EXCLUDED."end_date"
      )
    `,
		rows: rows,
	})
	if err != nil {
		return &updateInfoTerm{}, err
	}

	return &updateInfoTerm{
		count: result.RowsAffected(),
	}, nil
}
