ALTER TABLE "term"
  ADD COLUMN "created_at" timestamptz NOT NULL DEFAULT NOW()
, ADD COLUMN "updated_at" timestamptz NOT NULL DEFAULT NOW()
, ADD COLUMN "deleted_at" timestamptz;

CREATE TRIGGER "term_update_set_timestamp"
BEFORE UPDATE ON "term"
FOR EACH ROW EXECUTE FUNCTION update_set_timestamp();
