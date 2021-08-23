package db

import (
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	//"log"
)

type UpdateSummariesSections struct {
	Section         UpdateSummary
	SectionStaff    UpdateSummary
	SectionSchedule UpdateSummary
}

type sectionCurrentRow struct {
	id      pgtype.UUID
	section *data.CourseSection
}

type rowsSections struct {
	section         [][]interface{}
	sectionStaff    [][]interface{}
	sectionSchedule [][]interface{}
}

func rowsForSections(
	sections map[data.SectionKey]*data.CourseSection,
) rowsSections {
	rows := rowsSections{
		make([][]interface{}, 0, len(sections)),
		make([][]interface{}, 0, len(sections)),
		make([][]interface{}, 0, len(sections)),
	}

	for key, section := range sections {
		withKey := func(vals ...interface{}) []interface{} {
			return append([]interface{}{
				key.Course.Department,
				key.Course.Code,
				key.Course.Campus,
				key.Term,
				key.Section,
			}, vals...)
		}

		rows.section = append(rows.section, withKey(
			section.QuarterCredits,
			section.Status.String(),
			section.Seats.Enrolled,
			section.Seats.Capacity,
		))

		for _, staff := range section.Staff {
			rows.sectionStaff = append(rows.sectionStaff, withKey(staff))
		}

		for _, schedule := range section.Schedule {
			rows.sectionSchedule = append(rows.sectionSchedule, withKey(
				int(schedule.Days),
				schedule.Start.Std(),
				schedule.End.Std(),
				schedule.Location,
			))
		}
	}

	return rows
}

func (rows rowsSections) pushSection(h *handle) (UpdateSummary, error) {
	if _, err := h.tx.Exec(h.ctx, `
    CREATE TEMPORARY TABLE "section_tmp"
      ( "course_department" text    NOT NULL
      , "course_code"       text    NOT NULL
      , "course_campus"     text    NOT NULL
      , "term_code"         text    NOT NULL
      , "section"           integer NOT NULL
      , "quarter_credits"   integer NOT NULL
      , "status"            status  NOT NULL
      , "seats_enrolled"    integer NOT NULL
      , "seats_capacity"    integer NOT NULL
      )
    ON COMMIT DROP;
  `); err != nil {
		return UpdateSummary{}, fmt.Errorf("failed to create `section_tmp` table: %w", err)
	}

	countCopy, err := h.tx.CopyFrom(h.ctx, pgx.Identifier{"section_tmp"},
		[]string{
			"course_department",
			"course_code",
			"course_campus",
			"term_code",
			"section",
			"quarter_credits",
			"status",
			"seats_enrolled",
			"seats_capacity",
		},
		pgx.CopyFromRows(rows.section),
	)
	if err != nil {
		return UpdateSummary{}, fmt.Errorf("failed to populate `section_tmp`: %w", err)
	}

	tagUpsert, err := h.tx.Exec(h.ctx, `
    INSERT INTO "section"
      ( "course_id"
      , "term_code"
      , "section"
      , "quarter_credits"
      , "status"
      , "seats_enrolled"
      , "seats_capacity"
      )
    SELECT
        "course"."id"
      , "section_tmp"."term_code"
      , "section_tmp"."section"
      , "section_tmp"."quarter_credits"
      , "section_tmp"."status"
      , "section_tmp"."seats_enrolled"
      , "section_tmp"."seats_capacity"
    FROM "section_tmp"
    JOIN "course" ON
      "course"."department" = "section_tmp"."course_department"
      AND "course"."code" = "section_tmp"."course_code"
      AND "course"."campus" = "section_tmp"."course_campus"
    ON CONFLICT ("course_id", "term_code", "section")
    DO UPDATE SET
        "quarter_credits" = EXCLUDED."quarter_credits"
      , "status" = EXCLUDED."status"
      , "seats_enrolled" = EXCLUDED."seats_enrolled"
      , "seats_capacity" = EXCLUDED."seats_capacity"
      , "deleted_at" = NULL
    WHERE
      "section"."quarter_credits" <> EXCLUDED."quarter_credits"
      OR "section"."status" <> EXCLUDED."status"
      OR "section"."seats_enrolled" <> EXCLUDED."seats_enrolled"
      OR "section"."seats_capacity" <> EXCLUDED."seats_capacity"
      OR "section"."deleted_at" IS NOT NULL
  `)
	if err != nil {
		return UpdateSummary{}, fmt.Errorf("failed to upsert section entries: %w", err)
	}

	tagDelete, err := h.tx.Exec(h.ctx, `
    UPDATE "section" SET "deleted_at" = NOW()
    WHERE
      NOT EXISTS (
        SELECT NULL
        FROM "section_tmp"
        JOIN "course" ON
          "course"."department" = "section_tmp"."course_department"
          AND "course"."code" = "section_tmp"."course_code"
          AND "course"."campus" = "section_tmp"."course_campus"
        WHERE
          "section"."course_id" = "course"."id"
          AND "section"."term_code" = "section_tmp"."term_code"
          AND "section"."section" = "section_tmp"."section"
      )
      AND "deleted_at" IS NULL
  `)
	if err != nil {
		return UpdateSummary{}, fmt.Errorf("failed to delete old section entries: %w", err)
	}

	return UpdateSummary{
		Copied:   countCopy,
		Upserted: tagUpsert.RowsAffected(),
		Deleted:  tagDelete.RowsAffected(),
	}, nil

}

func (rows rowsSections) pushSectionStaff(h *handle) (UpdateSummary, error) {
	if _, err := h.tx.Exec(h.ctx, `
    CREATE TEMPORARY TABLE "section_staff_tmp"
      ( "course_department" text NOT NULL
      , "course_code"       text NOT NULL
      , "course_campus"     text NOT NULL
      , "section_term_code" text NOT NULL
      , "section_section"   integer NOT NULL
      , "staff_lingk_id"    text NOT NULL
      , PRIMARY KEY
          ( "course_department"
          , "course_code"
          , "course_campus"
          , "section_term_code"
          , "section_section"
          , "staff_lingk_id"
          )
      )
    ON COMMIT DROP;
  `); err != nil {
		return UpdateSummary{}, fmt.Errorf("failed to create section_staff_tmp: %w", err)
	}

	countCopy, err := h.tx.CopyFrom(h.ctx, pgx.Identifier{"section_staff_tmp"},
		[]string{
			"course_department",
			"course_code",
			"course_campus",
			"section_term_code",
			"section_section",
			"staff_lingk_id",
		},
		pgx.CopyFromRows(rows.sectionStaff),
	)
	if err != nil {
		return UpdateSummary{}, fmt.Errorf("failed to populate section_staff_tmp: %w", err)
	}

	tagUpsert, err := h.tx.Exec(h.ctx, `
    INSERT INTO "section_staff" ("section_id", "staff_id")
    SELECT "section"."id", "staff"."id"
    FROM "section_staff_tmp" AS "tmp"
    JOIN "course" ON
      ("course"."department", "course"."code", "course"."campus")
      = ("tmp"."course_department", "tmp"."course_code", "tmp"."course_campus")
    JOIN "section" ON
      ("section"."course_id", "section"."term_code", "section"."section")
      = ("course"."id", "tmp"."section_term_code", "tmp"."section_section")
    JOIN "staff" ON
      "staff"."lingk_id" = "tmp"."staff_lingk_id"
    ON CONFLICT ("section_id", "staff_id")
    DO UPDATE SET "deleted_at" = NULL
    WHERE "section_staff"."deleted_at" IS NOT NULL
  `)
	if err != nil {
		return UpdateSummary{}, fmt.Errorf("failed to upsert section_staff rows: %w", err)
	}

	tagDelete, err := h.tx.Exec(h.ctx, `
    UPDATE "section_staff" AS "current" SET "deleted_at" = NOW()
    WHERE
      NOT EXISTS
        ( SELECT NULL FROM "section_staff_tmp" AS "tmp"
          JOIN "course" ON
            ("course"."department", "course"."code", "course"."campus")
            = ("tmp"."course_department", "tmp"."course_code", "tmp"."course_campus")
          JOIN "section" ON
            ("section"."course_id", "section"."term_code", "section"."section")
            = ("course"."id", "tmp"."section_term_code", "tmp"."section_section")
          JOIN "staff" ON
            "staff"."lingk_id" = "tmp"."staff_lingk_id"
          WHERE
            ("current"."section_id", "current"."staff_id")
            = ("section"."id", "staff"."id")
        )
      AND "deleted_at" IS NULL
  `)
	if err != nil {
		return UpdateSummary{}, fmt.Errorf("failed to update deleted section_staff rows: %w", err)
	}

	return UpdateSummary{
		Copied:   countCopy,
		Upserted: tagUpsert.RowsAffected(),
		Deleted:  tagDelete.RowsAffected(),
	}, nil
}

func (rows rowsSections) pushSectionSchedule(h *handle) (UpdateSummary, error) {
	if _, err := h.tx.Exec(h.ctx, `
    CREATE TEMPORARY TABLE "section_schedule_tmp"
      ( "course_department" text NOT NULL
      , "course_code"       text NOT NULL
      , "course_campus"     text NOT NULL
      , "section_term_code" text NOT NULL
      , "section_section"   integer NOT NULL
      , "days"              smallint NOT NULL
      , "time_start"        time NOT NULL
      , "time_end"          time NOT NULL
      , "location"          text NOT NULL
      , PRIMARY KEY
          ( "course_department"
          , "course_code"
          , "course_campus"
          , "section_term_code"
          , "section_section"
          , "days"
          , "time_start"
          , "time_end"
          , "location"
          )
      )
    ON COMMIT DROP;
  `); err != nil {
		return UpdateSummary{}, fmt.Errorf("failed to create section_schedule_tmp: %w", err)
	}

	countCopy, err := h.tx.CopyFrom(h.ctx, pgx.Identifier{"section_schedule_tmp"},
		[]string{
			"course_department",
			"course_code",
			"course_campus",
			"section_term_code",
			"section_section",
			"days",
			"time_start",
			"time_end",
			"location",
		},
		pgx.CopyFromRows(rows.sectionSchedule),
	)
	if err != nil {
		return UpdateSummary{}, fmt.Errorf("failed to populate section_schedule_tmp: %w", err)
	}

	tagUpsert, err := h.tx.Exec(h.ctx, `
    INSERT INTO "section_schedule" AS "current"
      ( "section_id"
      , "days"
      , "time_start"
      , "time_end"
      , "location"
      )
    SELECT
      "section"."id"
    , "tmp"."days"
    , "tmp"."time_start"
    , "tmp"."time_end"
    , "tmp"."location"
    FROM "section_schedule_tmp" AS "tmp"
    JOIN "course" ON
      ("course"."department", "course"."code", "course"."campus")
      = ("tmp"."course_department", "tmp"."course_code", "tmp"."course_campus")
    JOIN "section" ON
      ("section"."course_id", "section"."term_code", "section"."section")
      = ("course"."id", "tmp"."section_term_code", "tmp"."section_section")
    ON CONFLICT ("section_id", "days", "time_start", "time_end", "location")
    DO UPDATE SET
      "location" = EXCLUDED."location"
    , "deleted_at" = NULL
    WHERE
      "current"."location" <> EXCLUDED."location"
      OR "current"."deleted_at" IS NOT NULL
  `)
	if err != nil {
		return UpdateSummary{}, fmt.Errorf("failed to upsert section_schedule rows: %w", err)
	}

	tagDelete, err := h.tx.Exec(h.ctx, `
    UPDATE "section_schedule" AS "current" SET "deleted_at" = NOW()
    WHERE
      NOT EXISTS
        ( SELECT NULL FROM "section_schedule_tmp" AS "tmp"
          JOIN "course" ON
            ("course"."department", "course"."code", "course"."campus")
            = ("tmp"."course_department", "tmp"."course_code", "tmp"."course_campus")
          JOIN "section" ON
            ("section"."course_id", "section"."term_code", "section"."section")
            = ("course"."id", "tmp"."section_term_code", "tmp"."section_section")
          WHERE
            ( "section"."id"
            , "tmp"."days"
            , "tmp"."time_start"
            , "tmp"."time_end"
            , "tmp"."location"
            )
            =
            ( "current"."section_id"
            , "current"."days"
            , "current"."time_start"
            , "current"."time_end"
            , "current"."location"
            )
        )
      AND "deleted_at" IS NULL
  `)
	if err != nil {
		return UpdateSummary{}, fmt.Errorf("failed to update deleted section_schedule rows: %w", err)
	}

	return UpdateSummary{
		Copied:   countCopy,
		Upserted: tagUpsert.RowsAffected(),
		Deleted:  tagDelete.RowsAffected(),
	}, nil
}

func (h *handle) updateSections(
	sections map[data.SectionKey]*data.CourseSection,
) (UpdateSummariesSections, error) {

	rows := rowsForSections(sections)

	summarySection, err := rows.pushSection(h)
	if err != nil {
		return UpdateSummariesSections{}, err
	}

	summarySectionStaff, err := rows.pushSectionStaff(h)
	if err != nil {
		return UpdateSummariesSections{}, err
	}

	summarySectionSchedule, err := rows.pushSectionSchedule(h)
	if err != nil {
		return UpdateSummariesSections{}, err
	}

	return UpdateSummariesSections{
		Section:         summarySection,
		SectionStaff:    summarySectionStaff,
		SectionSchedule: summarySectionSchedule,
	}, nil

}
