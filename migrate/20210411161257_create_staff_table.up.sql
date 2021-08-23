CREATE TABLE "staff"
( "id"         uuid PRIMARY KEY DEFAULT gen_random_uuid()
, "lingk_id"   text NOT NULL UNIQUE
, "first_name" text NOT NULL
, "last_name"  text NOT NULL
, "created_at" timestamptz NOT NULL DEFAULT NOW()
, "updated_at" timestamptz NOT NULL DEFAULT NOW()
);

CREATE TRIGGER "staff_update_set_timestamp"
BEFORE UPDATE ON "staff"
FOR EACH ROW EXECUTE FUNCTION update_set_timestamp();
