INSERT INTO "course"
( "department"
, "code"
, "campus"
, "name"
, "description"
)
SELECT
  "department"
, "code"
, "campus"
, "name"
, "description"
FROM "course_tmp"
ON CONFLICT ( "department", "code", "campus" ) DO UPDATE SET
  "name"        = EXCLUDED."name"
, "description" = EXCLUDED."description"
, "deleted_at"  = NULL
WHERE
  "course"."name" <> EXCLUDED."name"
  OR "course"."description" <> EXCLUDED."description"
  OR "course"."deleted_at" IS NOT NULL;
