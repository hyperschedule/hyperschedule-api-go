UPDATE "term" AS "current" SET "deleted_at" = NOW()
WHERE
  NOT EXISTS (
    SELECT NULL FROM "term_tmp" AS "tmp"
    WHERE "tmp"."code" = "current"."code"
  )
  AND "current"."deleted_at" IS NULL;
