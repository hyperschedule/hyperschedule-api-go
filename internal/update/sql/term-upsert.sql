INSERT INTO "term"
( "code", "semester", "date_start", "date_end" )
SELECT
  "code", "semester", "date_start", "date_end"
FROM "term_tmp"
ON CONFLICT ("code") DO UPDATE SET
  "semester"   = EXCLUDED."semester"
, "date_start" = EXCLUDED."date_start"
, "date_end"   = EXCLUDED."date_end"
WHERE
  ( "term"."semester", "term"."date_start", "term"."date_end" )
  <> ( EXCLUDED."semester", EXCLUDED."date_start", EXCLUDED."date_end" );
