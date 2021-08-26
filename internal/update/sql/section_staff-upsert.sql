INSERT INTO "section_staff" ("section_id", "staff_id")
SELECT "section"."id", "staff"."id"
FROM "section_staff_tmp" AS "tmp"
JOIN "course" ON
  ("course"."department", "course"."code", "course"."campus")
  = ("tmp"."course_department", "tmp"."course_code", "tmp"."course_campus")
JOIN "section" ON
  ("section"."course_id", "section"."term_code", "section"."section")
  = ("course"."id", "tmp"."term_code", "tmp"."section")
JOIN "staff" ON
  "staff"."lingk_id" = "tmp"."staff_lingk_id"
ON CONFLICT ("section_id", "staff_id")
DO UPDATE SET "deleted_at" = NULL
WHERE "section_staff"."deleted_at" IS NOT NULL

