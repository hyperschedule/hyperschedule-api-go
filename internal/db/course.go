package db

import (
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
	"github.com/jackc/pgx/v4"
)

type updateInfoCourse struct {
	Upserted int64
	Deleted  int64
}

func (h *handle) updateCourses(courses map[data.CourseKey]*data.Course) (*updateInfoCourse, error) {
	rows := make([][]interface{}, 0, len(courses))
	for key, course := range courses {
		rows = append(rows, []interface{}{
			key.Department,
			key.Code,
			key.Campus,
			course.Name,
			course.Description,
		})
	}

	if _, err := h.tx.Exec(h.ctx, `
    CREATE TEMPORARY TABLE "course_tmp"
    ( "department"  text NOT NULL
    , "code"        text NOT NULL
    , "campus"      text NOT NULL
    , "name"        text NOT NULL
    , "description" text NOT NULL
    , UNIQUE ("department", "code", "campus")
    )
    ON COMMIT DROP;
  `); err != nil {
		return nil, fmt.Errorf("failed to create course_tmp table: %w", err)
	}

	if _, err := h.tx.CopyFrom(h.ctx, pgx.Identifier{"course_tmp"},
		[]string{"department", "code", "campus", "name", "description"},
		pgx.CopyFromRows(rows),
	); err != nil {
		return nil, fmt.Errorf("failed to populate course_tmp table: %w", err)
	}

	tagUpsert, err := h.tx.Exec(h.ctx, `
    INSERT INTO "course"
      ( "department"
      , "code"
      , "campus"
      , "name"
      , "description"
      )
    SELECT
        "department"
      , "code"
      , "campus"
      , "name"
      , "description"
    FROM "course_tmp"
    ON CONFLICT ("department", "code", "campus") DO UPDATE SET
      "name" = EXCLUDED."name"
    , "description" = EXCLUDED."description"
    , "deleted_at" = NULL
    WHERE
      "course"."name" <> EXCLUDED."name"
      OR "course"."description" <> EXCLUDED."description"
      OR "course"."deleted_at" IS NOT NULL;
  `)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert course entries: %w", err)
	}

	tagDelete, err := h.tx.Exec(h.ctx, `
    UPDATE "course" SET "deleted_at" = NOW()
    WHERE
      NOT EXISTS (
        SELECT NULL FROM "course_tmp" WHERE
          "course"."department" = "course_tmp"."department"
          AND "course"."code" = "course_tmp"."code"
          AND "course"."campus" = "course_tmp"."campus"
      )
      AND "deleted_at" IS NULL
  `)
	if err != nil {
		return nil, fmt.Errorf("failed to delete old course entries: %w", err)
	}

	return &updateInfoCourse{
		Upserted: tagUpsert.RowsAffected(),
		Deleted:  tagDelete.RowsAffected(),
	}, nil

}
