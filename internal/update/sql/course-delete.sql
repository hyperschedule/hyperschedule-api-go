UPDATE "course" SET "deleted_at" = NOW()
WHERE
  NOT EXISTS (
    SELECT NULL FROM "course_tmp" WHERE
      "course"."department" = "course_tmp"."department"
      AND "course"."code" = "course_tmp"."code"
      AND "course"."campus" = "course_tmp"."campus"
  )
  AND "deleted_at" IS NULL
