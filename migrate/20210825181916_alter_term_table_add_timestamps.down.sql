ALTER TABLE "term"
  DROP COLUMN "created_at"
, DROP COLUMN "updated_at"
, DROP COLUMN "deleted_at";

DROP TRIGGER "term_update_set_timestamp";
