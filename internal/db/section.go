package db

import (
	"errors"
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
	"github.com/jackc/pgtype"
	"time"
)

type updateInfoSection struct {
	IndexNew int64
	Delete   int
}

type sectionCurrentRow struct {
	id      pgtype.UUID
	section *data.CourseSection
}

func updateSections(
	h *handle,
	sections map[data.SectionKey]*data.CourseSection,
) (*updateInfoSection, error) {

	indexNew, err := updateSectionPopulateIndex(h, sections)
	if err != nil {
		return nil, fmt.Errorf("index failed: %w", err)
	}

	currentRows, err := updateSectionFetchCurrent(h)
	if err != nil {
		return nil, fmt.Errorf("current: %w", err)
	}

	delCount, err := updateSectionDel(h, currentRows, sections)
	if err != nil {
		return nil, fmt.Errorf("del: %w", err)
	}

	////updateSectionSnapshots

	statusRows, err := updateSectionStatus(h, currentRows, sections)
	if err != nil {
		return nil, fmt.Errorf("status: %w", err)
	}

	_ = statusRows

	fmt.Printf("UPDATE SECTION\n")

	//updateSectionLatest

	return &updateInfoSection{
		IndexNew: indexNew,
		Delete:   delCount,
	}, nil
}

func updateSectionPopulateIndex(h *handle,
	sections map[data.SectionKey]*data.CourseSection,
) (int64, error) {
	rows := make([][]interface{}, 0)
	for key := range sections {
		rows = append(rows, []interface{}{
			key.Course.Department,
			key.Course.Code,
			key.Course.Campus,
			key.Term,
			key.Section,
		})
	}
	result, err := batchInsert(h, &batchOpts{
		tmpName: "section_tmp",
		tmpCols: []string{"department", "code", "campus", "term_code", "section"},
		sqlTmpDefs: `
      "department" text NOT NULL
    , "code" text NOT NULL
    , "campus" text NOT NULL
    , "term_code" text NOT NULL
    , "section" int NOT NULL
    , PRIMARY KEY ("department", "code", "campus", "term_code", "section")
    `,
		sqlInsertSelect: `
      INSERT INTO "section" ("course_id", "term_code", "section")
      SELECT
        "course"."id"
      , "tmp"."term_code"
      , "tmp"."section"
      FROM "section_tmp" AS "tmp"
      JOIN "course" ON
        "course"."department" = "tmp"."department"
        AND "course"."code" = "tmp"."code"
        AND "course"."campus" = "tmp"."campus"
      ON CONFLICT ("course_id", "term_code", "section") DO NOTHING
    `,
		rows: rows,
	})
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func updateSectionFetchCurrent(
	h *handle,
) (map[data.SectionKey]*sectionCurrentRow, error) {
	rows, err := h.tx.Query(h.ctx, `
    SELECT
      "section"."id"
    , "section"."term_code"
    , "section"."section"
    , "course"."department"
    , "course"."code"
    , "course"."campus"
    , "snapshot"."quarter_credits"
    , "status"."status"
    , "status"."enrolled"
    , "status"."capacity"
    FROM "section_latest" AS "latest"
    JOIN "section" ON "section"."id" = "latest"."section_id"
    JOIN "course" ON "course"."id" = "section"."course_id"
    JOIN "section_snapshot" AS "snapshot" ON 
      "snapshot"."section_id" = "section"."id"
      AND "snapshot"."time" = "latest"."snapshot_time"
    JOIN "section_status" AS "status" ON
      "status"."section_id" = "section"."id"
      AND "status"."time" = "latest"."status_time"
  `)
	if err != nil {
		return nil, err
	}

	sections := make(map[data.SectionKey]*sectionCurrentRow)
	ids := make([]pgtype.UUID, 0)
	for rows.Next() {
		key := data.SectionKey{}
		row := &sectionCurrentRow{section: &data.CourseSection{}}
		var statusStr string
		if err := rows.Scan(
			&row.id,
			&key.Term,
			&key.Section,
			&key.Course.Department,
			&key.Course.Code,
			&key.Course.Campus,
			&row.section.QuarterCredits,
			&statusStr,
			&row.section.Seats.Enrolled,
			&row.section.Seats.Capacity,
		); err != nil {
			return nil, err
		}
		status, err := data.StatusFromString(statusStr)
		if err != nil {
			return nil, err
		}
		row.section.Status = status
		ids = append(ids, row.id)
		sections[key] = row
	}

	if err := updateSectionFetchCurrentSchedule(h, sections); err != nil {
		return nil, err
	}

	if err := updateSectionFetchCurrentStaff(h, sections); err != nil {
		return nil, err
	}

	return sections, nil
}

func updateSectionFetchCurrentSchedule(h *handle, sections map[data.SectionKey]*sectionCurrentRow) error {
	rows, err := h.tx.Query(h.ctx, `
    SELECT
      "section"."term_code"
    , "section"."section"
    , "course"."department"
    , "course"."code"
    , "course"."campus"
    , "schedule"."days"
    , "schedule"."location"
    , "schedule"."start_time"
    , "schedule"."end_time"
    FROM "section_latest" AS "latest"
    JOIN "section" ON "section"."id" = "latest"."section_id"
    JOIN "course" ON "course"."id" = "section"."course_id"
    JOIN "section_schedule" AS "schedule" ON
      "schedule"."section_id" = "section"."id"
      AND "schedule"."snapshot_time" = "latest"."snapshot_time"
  `)
	if err != nil {
		return err
	}

	for rows.Next() {
		key := data.SectionKey{}
		schedule := &data.Schedule{}
		var startTime, endTime time.Time
		if err := rows.Scan(
			&key.Term,
			&key.Section,
			&key.Course.Department,
			&key.Course.Code,
			&key.Course.Campus,
			&schedule.Days,
			&schedule.Location,
			&startTime,
			&endTime,
		); err != nil {
			return err
		}
		schedule.Start = data.TimeFromStd(startTime)
		schedule.End = data.TimeFromStd(endTime)

		row, ok := sections[key]
		if !ok {
			return errors.New("missing section key (unexpected)")
		}
		row.section.Schedule = append(row.section.Schedule, schedule)
	}
	return nil
}

func updateSectionFetchCurrentStaff(h *handle, sections map[data.SectionKey]*sectionCurrentRow) error {
	rows, err := h.tx.Query(h.ctx, `
    SELECT
      "section"."term_code"
    , "section"."section"
    , "course"."department"
    , "course"."code"
    , "course"."campus"
    , "section_staff"."staff_id"
    FROM "section_latest" AS "latest"
    JOIN "section" ON "section"."id" = "latest"."section_id"
    JOIN "course" ON "course"."id" = "section"."course_id"
    JOIN "section_staff" ON
      "section_staff"."section_id" = "section"."id"
      AND "section_staff"."snapshot_time" = "latest"."snapshot_time"
  `)
	if err != nil {
		return err
	}

	for rows.Next() {
		key := data.SectionKey{}
		var staffId string
		if err := rows.Scan(
			&key.Term,
			&key.Section,
			&key.Course.Department,
			&key.Course.Code,
			&key.Course.Campus,
			&staffId,
		); err != nil {
			return err
		}

		row, ok := sections[key]
		if !ok {
			return errors.New("missing section :(")
		}
		row.section.Staff = append(row.section.Staff, staffId)
	}

	return nil
}

func updateSectionDel(h *handle,
	current map[data.SectionKey]*sectionCurrentRow,
	next map[data.SectionKey]*data.CourseSection,
) (int, error) {
	delIds := make([]pgtype.UUID, 0)
	delRows := make([][]interface{}, 0)
	for key, row := range current {
		if _, ok := next[key]; !ok {
			delIds = append(delIds, row.id)
			delRows = append(delRows, []interface{}{row.id})
		}
	}

	if _, err := batchInsert(h, &batchOpts{
		tmpName: "section_delete_tmp",
		tmpCols: []string{"section_id"},
		sqlTmpDefs: `
      "section_id" uuid NOT NULL
    `,
		sqlInsertSelect: `
      INSERT INTO "section_delete" ("section_id")
      SELECT "section_id" FROM "section_delete_tmp"
    `,
		rows: delRows,
	}); err != nil {
		return 0, err
	}

	result, err := h.tx.Exec(h.ctx, `
    DELETE FROM "section_latest"
    WHERE "section_id" = ANY($1) 
  `, delIds)
	if err != nil {
		return 0, err
	}

	if result.RowsAffected() != int64(len(delIds)) {
		return 0, errors.New("mismatch section delete counts")
	}

	return len(delIds), nil
}

type idTime struct {
	id   pgtype.UUID
	time time.Time
}

func updateSectionStatus(
	h *handle,
	current map[data.SectionKey]*sectionCurrentRow,
	next map[data.SectionKey]*data.CourseSection,
) ([]idTime, error) {
	rows := make([][]interface{}, 0)
	for key, section := range next {
		old, ok := current[key]
		if !ok || section.Status != old.section.Status || section.Seats != old.section.Seats {
			rows = append(rows, []interface{}{
				key.Term,
				key.Section,
				key.Course.Department,
				key.Course.Code,
				key.Course.Campus,
				section.Status.String(),
				section.Seats.Enrolled,
				section.Seats.Capacity,
			})
		}
	}

	result, err := batchInsertReturns(h, &batchOpts{
		tmpName: "section_status_tmp",
		tmpCols: []string{
			"term_code",
			"section",
			"department",
			"code",
			"campus",
			"status",
			"enrolled",
			"capacity",
		},
		sqlTmpDefs: `
      "term_code" text NOT NULL
    , "section" int NOT NULL
    , "department" text NOT NULL
    , "code" text NOT NULL
    , "campus" text NOT NULL
    , "status" status NOT NULL
    , "enrolled" int NOT NULL
    , "capacity" int NOT NULL
    , PRIMARY KEY ("term_code", "section", "department", "code", "campus")
    `,
		sqlInsertSelect: `
      INSERT INTO "section_status" ("section_id", "status", "enrolled", "capacity")
      SELECT 
        "section"."id"
      , "tmp"."status"
      , "tmp"."enrolled"
      , "tmp"."capacity"
      FROM "section_status_tmp" AS "tmp"
      JOIN "section" ON
        "section"."term_code" = "tmp"."term_code"
        AND "section"."section" = "tmp"."section"
      JOIN "course" ON
        "course"."id" = "section"."course_id"
        AND "course"."department" = "tmp"."department"
        AND "course"."code" = "tmp"."code"
        AND "course"."campus" = "tmp"."campus"
      RETURNING "section_id", "time"
    `,
		rows: rows,
	})
	if err != nil {
		return nil, err
	}

	pts := make([]idTime, len(rows))
	for result.Next() {
		pt := idTime{}
		if err := result.Scan(&pt.id, &pt.time); err != nil {
			return nil, err
		}
		pts = append(pts, pt)
	}

	return pts, nil
}
