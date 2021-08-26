UPDATE "section_schedule" AS "current" SET "deleted_at" = NOW()
WHERE
  NOT EXISTS
    ( SELECT NULL FROM "section_schedule_tmp" AS "tmp"
      JOIN "course" ON
        ("course"."department", "course"."code", "course"."campus")
        = ("tmp"."course_department", "tmp"."course_code", "tmp"."course_campus")
      JOIN "section" ON
        ("section"."course_id", "section"."term_code", "section"."section")
        = ("course"."id", "tmp"."term_code", "tmp"."section")
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

