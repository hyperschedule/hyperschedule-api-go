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
  = ("course"."id", "tmp"."term_code", "tmp"."section")
ON CONFLICT ("section_id", "days", "time_start", "time_end", "location")
DO UPDATE SET
  "location" = EXCLUDED."location"
, "deleted_at" = NULL
WHERE
  "current"."location" <> EXCLUDED."location"
  OR "current"."deleted_at" IS NOT NULL

