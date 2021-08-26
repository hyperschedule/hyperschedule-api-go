UPDATE "staff" AS "current" SET "deleted_at" = NOW()
WHERE
  NOT EXISTS (
    SELECT NULL FROM "staff_tmp" AS "tmp"
    WHERE "tmp"."lingk_id" = "current"."lingk_id"
  )
  AND "current"."deleted_at" IS NULL;
