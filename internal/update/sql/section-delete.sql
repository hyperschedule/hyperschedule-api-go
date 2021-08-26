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

