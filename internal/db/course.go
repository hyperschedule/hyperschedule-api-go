package db

import (
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
	"github.com/jackc/pgtype"
	"time"
)

type updateInfoCourse struct {
	Updated int
	Del     *updateInfoCourseDel
}

type updateInfoCourseDel struct {
	Expect int
	Actual int64
}

func updateCourses(h *handle, courses map[data.CourseKey]*data.Course) (*updateInfoCourse, error) {
	if err := populateCourseIndex(h, courses); err != nil {
		return nil, err
	}

	prev, err := fetchCurrentRows(h)
	if err != nil {
		return nil, err
	}

	delInfo, err := courseProcessDeletedRows(h, prev, courses)
	if err != nil {
		return nil, err
	}

	// insert new/updated course snapshots and update latest rows
	snapshotRows := make([][]interface{}, 0)
	for key, course := range courses {
		last, ok := prev[key]
		if !ok || *course != *last.course {
			snapshotRows = append(snapshotRows, []interface{}{
				key.Department,
				key.Code,
				key.Campus,
				course.Name,
				course.Description,
			})
		}
	}
	snapshotIdRows, err := batchInsertReturns(h, &batchOpts{
		tmpName: "course_snapshot_tmp",
		tmpCols: []string{"department", "code", "campus", "name", "description"},
		sqlTmpDefs: `
      "department" text NOT NULL
    , "code" text NOT NULL
    , "campus" text NOT NULL
    , "name" text NOT NULL
    , "description" text NOT NULL
    , PRIMARY KEY ("department", "code", "campus")
    `,
		sqlInsertSelect: `
      INSERT INTO "course_snapshot" 
        ("course_id", "name", "description")
      SELECT "course"."id", "tmp"."name", "tmp"."description" 
      FROM "course_snapshot_tmp" AS "tmp"
      JOIN "course" ON 
        "course"."department" = "tmp"."department"
        AND "course"."code" = "tmp"."code"
        AND "course"."campus" = "tmp"."campus"
      RETURNING "course_id", "time"
    `,
		rows: snapshotRows,
	})
	if err != nil {
		return nil, err
	}
	latestRows := make([][]interface{}, 0)
	for snapshotIdRows.Next() {
		var id pgtype.UUID
		var t time.Time
		if err := snapshotIdRows.Scan(&id, &t); err != nil {
			return nil, err
		}
		latestRows = append(latestRows, []interface{}{id, t})
	}

	if _, err := batchInsert(h, &batchOpts{
		tmpName: "course_latest_tmp",
		tmpCols: []string{"course_id", "snapshot_time"},
		sqlTmpDefs: `
      "course_id" uuid PRIMARY KEY
    , "snapshot_time" timestamptz NOT NULL
    `,
		sqlInsertSelect: `
      INSERT INTO "course_latest" ("course_id", "snapshot_time")
      SELECT "course_id", "snapshot_time" FROM "course_latest_tmp"
      ON CONFLICT ("course_id") DO UPDATE SET 
        "snapshot_time" = EXCLUDED."snapshot_time"
    `,
		rows: latestRows,
	}); err != nil {
		return nil, err
	}

	return &updateInfoCourse{
		Updated: len(snapshotRows),
		Del:     delInfo,
	}, nil
}

// insert/create course index entries to ensure each course key (department,
// code, campus) has a UUID
func populateCourseIndex(h *handle, courses map[data.CourseKey]*data.Course) error {
	rows := [][]interface{}{}
	for key := range courses {
		rows = append(rows, []interface{}{key.Department, key.Code, key.Campus})
	}
	if _, err := batchInsert(h, &batchOpts{
		tmpName: "course_tmp",
		tmpCols: []string{"department", "code", "campus"},
		sqlTmpDefs: `
      "department" text NOT NULL
    , "code" text NOT NULL
    , "campus" text NOT NULL
    , PRIMARY KEY ("department", "code", "campus")
    `,
		sqlInsertSelect: `
      INSERT INTO "course" ("department", "code", "campus")
      SELECT "department", "code", "campus" FROM "course_tmp"
      ON CONFLICT ("department", "code", "campus") DO NOTHING
    `,
		rows: rows,
	}); err != nil {
		return err
	}
	return nil
}

type dbCourseRow struct {
	id     pgtype.UUID
	course *data.Course
}

// fetch last snapshot rows
func fetchCurrentRows(h *handle) (map[data.CourseKey]*dbCourseRow, error) {
	lastRows, err := h.tx.Query(h.ctx, `
    SELECT
      "course"."id"
    , "course"."department"
    , "course"."code"
    , "course"."campus"
    , "snap"."name"
    , "snap"."description"
    FROM "course_latest" AS "latest"
    JOIN "course" ON 
      "course"."id" = "latest"."course_id"
    JOIN "course_snapshot" AS "snap" ON 
      "snap"."course_id" = "course"."id"
      AND "snap"."time" = "latest"."snapshot_time"
  `)
	if err != nil {
		return nil, err
	}

	courses := make(map[data.CourseKey]*dbCourseRow, 0)
	for lastRows.Next() {
		key := data.CourseKey{}
		row := &dbCourseRow{course: &data.Course{}}
		if err := lastRows.Scan(
			&row.id,
			&key.Department,
			&key.Code,
			&key.Campus,
			&row.course.Name,
			&row.course.Description,
		); err != nil {
			return nil, err
		}
		courses[key] = row
	}

	return courses, nil
}

func courseProcessDeletedRows(
	h *handle,
	prev map[data.CourseKey]*dbCourseRow,
	courses map[data.CourseKey]*data.Course,
) (*updateInfoCourseDel, error) {
	// compute diff: delete rows from latest table & add delete entries
	delIds := make([]pgtype.UUID, 0)
	delRows := make([][]interface{}, 0)
	for key, row := range prev {
		if _, ok := courses[key]; !ok {
			delIds = append(delIds, row.id)
			delRows = append(delRows, []interface{}{row.id})
		}
	}

	if _, err := batchInsert(h, &batchOpts{
		tmpName:    "course_delete_tmp",
		tmpCols:    []string{"course_id"},
		sqlTmpDefs: `"course_id" uuid PRIMARY KEY`,
		sqlInsertSelect: `
      INSERT INTO "course_delete" ("course_id")
      SELECT "course_id" FROM "course_delete_tmp"
    `,
		rows: delRows,
	}); err != nil {
		return nil, err
	}

	delResult, err := h.tx.Exec(h.ctx,
		`DELETE FROM "course_latest" WHERE "course_id" = ANY($1)`,
		delIds,
	)
	if err != nil {
		return nil, err
	}

	return &updateInfoCourseDel{
		Expect: len(delIds),
		Actual: delResult.RowsAffected(),
	}, nil
}
