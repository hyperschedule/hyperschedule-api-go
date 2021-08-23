CREATE TABLE "course"
( "id" uuid PRIMARY KEY DEFAULT gen_random_uuid()
, "department" text NOT NULL
, "code" text NOT NULL
, "campus" text NOT NULL
, "name" text NOT NULL
, "description" text NOT NULL
, "deleted_at" timestamptz
, "created_at" timestamptz NOT NULL DEFAULT NOW()
, "updated_at" timestamptz NOT NULL DEFAULT NOW()
, UNIQUE ("department", "code", "campus")
);

CREATE TRIGGER "course_update_set_timestamp"
BEFORE UPDATE ON "course"
FOR EACH ROW EXECUTE FUNCTION update_set_timestamp();
