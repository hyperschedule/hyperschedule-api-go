INSERT INTO "section"
  ( "course_id"
  , "term_code"
  , "section"
  , "quarter_credits"
  , "status"
  , "seats_enrolled"
  , "seats_capacity"
  , "perms"
  )
SELECT
    "course"."id"
  , "section_tmp"."term_code"
  , "section_tmp"."section"
  , "section_tmp"."quarter_credits"
  , "section_tmp"."status"
  , "section_tmp"."seats_enrolled"
  , "section_tmp"."seats_capacity"
  , "section_tmp"."perms"
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
  , "perms" = EXCLUDED."perms"
  , "deleted_at" = NULL
WHERE
  "section"."quarter_credits" <> EXCLUDED."quarter_credits"
  OR "section"."status" <> EXCLUDED."status"
  OR "section"."seats_enrolled" <> EXCLUDED."seats_enrolled"
  OR "section"."seats_capacity" <> EXCLUDED."seats_capacity"
  OR "section"."perms" <> EXCLUDED."perms"
  OR "section"."deleted_at" IS NOT NULL

