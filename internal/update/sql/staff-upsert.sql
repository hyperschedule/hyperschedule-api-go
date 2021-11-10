INSERT INTO "staff"
( "lingk_id", "name_first", "name_last", "alt" )
SELECT
  "lingk_id", "name_first", "name_last", "alt"
FROM "staff_tmp"
ON CONFLICT ( "lingk_id" ) DO UPDATE SET
  "name_first" = EXCLUDED."name_first"
, "name_last"  = EXCLUDED."name_last"
, "alt"        = EXCLUDED."alt"
, "deleted_at" = NULL
WHERE
  "staff"."name_first" <> EXCLUDED."name_first"
  OR "staff"."name_last" <> EXCLUDED."name_last"
  OR "staff"."alt" IS DISTINCT FROM EXCLUDED."alt"
  OR "staff"."deleted_at" IS NOT NULL;
