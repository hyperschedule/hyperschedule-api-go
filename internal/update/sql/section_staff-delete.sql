UPDATE "section_staff" AS "current" SET "deleted_at" = NOW()
WHERE
  NOT EXISTS
    ( SELECT NULL FROM "section_staff_tmp" AS "tmp"
      JOIN "course" ON
        ("course"."department", "course"."code", "course"."campus")
        = ("tmp"."course_department", "tmp"."course_code", "tmp"."course_campus")
      JOIN "section" ON
        ("section"."course_id", "section"."term_code", "section"."section")
        = ("course"."id", "tmp"."term_code", "tmp"."section")
      JOIN "staff" ON
        "staff"."lingk_id" = "tmp"."staff_lingk_id"
      WHERE
        ("current"."section_id", "current"."staff_id")
        = ("section"."id", "staff"."id")
    )
  AND "deleted_at" IS NULL

