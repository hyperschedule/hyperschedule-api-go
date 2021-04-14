CREATE TYPE "status" AS ENUM ( 'open', 'closed', 'reopened' );

CREATE TABLE "section" 
( "id" uuid PRIMARY KEY DEFAULT gen_random_uuid()
, "term_code" text NOT NULL
, "course_code" text NOT NULL
, "campus" text NOT NULL
, "section" int NOT NULL
, UNIQUE ("term_code", "course_code", "campus", "section")
, FOREIGN KEY ("course_code", "campus") 
  REFERENCES "course" ("code", "campus")
);

CREATE TABLE "section_snapshot"
( "section_id" uuid NOT NULL REFERENCES "section"
, "time" timestamptz NOT NULL DEFAULT now()
, "quarter_credits" int NOT NULL
, PRIMARY KEY ("section_id", "time")
);

CREATE TABLE "section_schedule"
( "section_id" uuid NOT NULL
, "snapshot_time" timestamptz NOT NULL
, "location"    text NOT NULL
, "days"        bit(7) NOT NULL
, "start_time"  time NOT NULL
, "end_time"    time NOT NULL
, FOREIGN KEY ("section_id", "snapshot_time")
  REFERENCES "section_snapshot"
);

CREATE TABLE "section_staff"
( "section_id" uuid NOT NULL
, "snapshot_time" timestamptz NOT NULL
, "staff_id" uuid NOT NULL REFERENCES "staff"
, PRIMARY KEY ("section_id", "snapshot_time", "staff_id")
, FOREIGN KEY ("section_id", "snapshot_time")
  REFERENCES "section_snapshot"
);

CREATE TABLE "section_status"
( "section_id" uuid NOT NULL REFERENCES "section"
, "time" timestamptz NOT NULL DEFAULT now()
, "status" status NOT NULL
, "enrolled" int NOT NULL
, "capacity" int NOT NULL
, PRIMARY KEY ("section_id", "time")
);

CREATE TABLE "section_latest"
( "section_id" uuid PRIMARY KEY REFERENCES "section"
, "snapshot_time" timestamptz NOT NULL
, "status_time" timestamptz NOT NULL
, FOREIGN KEY ("section_id", "snapshot_time")
  REFERENCES "section_snapshot"
, FOREIGN KEY ("section_id", "status_time")
  REFERENCES "section_status"
);

CREATE OR REPLACE FUNCTION section_immut()
RETURNS TRIGGER LANGUAGE 'plpgsql'
AS $$ BEGIN RAISE EXCEPTION 'section_immut'; END $$;

CREATE TRIGGER "section_immut" 
BEFORE UPDATE ON "section"
EXECUTE FUNCTION section_immut();

CREATE TRIGGER "section_snapshot_immut" 
BEFORE UPDATE ON "section_snapshot"
EXECUTE FUNCTION section_immut();

CREATE TRIGGER "section_schedule_immut" 
BEFORE UPDATE ON "section_schedule"
EXECUTE FUNCTION section_immut();

CREATE TRIGGER "section_staff_immut" 
BEFORE UPDATE ON "section_staff"
EXECUTE FUNCTION section_immut();

CREATE TRIGGER "section_status_immut" 
BEFORE UPDATE ON "section_status"
EXECUTE FUNCTION section_immut();
