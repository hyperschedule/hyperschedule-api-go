CREATE TYPE "status" AS ENUM ( 'open', 'closed', 'reopened' );

CREATE TABLE "section"
( "id" uuid PRIMARY KEY DEFAULT gen_random_uuid()
, "course_id" uuid NOT NULL REFERENCES "course"
, "term_code" text NOT NULL REFERENCES "term"
, "section" int NOT NULL
, "quarter_credits" int NOT NULL
, "status" status NOT NULL
, "seats_enrolled" int NOT NULL
, "seats_capacity" int NOT NULL
, "created_at" timestamptz NOT NULL DEFAULT NOW()
, "updated_at" timestamptz NOT NULL DEFAULT NOW()
, "deleted_at" timestamptz
, UNIQUE ("course_id", "term_code", "section")
);

CREATE TABLE "section_schedule"
( "section_id" uuid NOT NULL REFERENCES "section"
, "days"       smallint NOT NULL
, "time_start" time NOT NULL
, "time_end"   time NOT NULL
, "location"   text NOT NULL
, "created_at" timestamptz NOT NULL DEFAULT NOW()
, "updated_at" timestamptz NOT NULL DEFAULT NOW()
, "deleted_at" timestamptz
, PRIMARY KEY ("section_id", "days", "time_start", "time_end", "location")
);

CREATE TABLE "section_staff"
( "section_id" uuid NOT NULL REFERENCES "section"
, "staff_id"   uuid NOT NULL REFERENCES "staff"
, "created_at" timestamptz NOT NULL DEFAULT NOW()
, "updated_at" timestamptz NOT NULL DEFAULT NOW()
, "deleted_at" timestamptz
, PRIMARY KEY ("section_id", "staff_id")
);

CREATE TRIGGER "section_staff_update_set_timestamp"
BEFORE UPDATE ON "section_staff"
FOR EACH ROW EXECUTE FUNCTION update_set_timestamp();

CREATE TRIGGER "section_schedule_update_set_timestamp"
BEFORE UPDATE ON "section_schedule"
FOR EACH ROW EXECUTE FUNCTION update_set_timestamp();

CREATE TRIGGER "section_update_set_timestamp"
BEFORE UPDATE ON "section"
FOR EACH ROW EXECUTE FUNCTION update_set_timestamp();
