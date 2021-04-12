CREATE TYPE "status" AS ENUM ( 'open', 'closed', 'reopened' );

CREATE TABLE "section_snapshot"
( "at" timestamptz NOT NULL DEFAULT now()
, "term_code" text NOT NULL
, "course_code" text NOT NULL
, "campus" text NOT NULL
, "section" int NOT NULL
, "quarter_credits" int NOT NULL
, PRIMARY KEY ("at", "term_code", "course_code", "campus", "section")
, FOREIGN KEY ("course_code", "campus") 
  REFERENCES "course" ("code", "campus")
);

CREATE TABLE "section_schedule"
( "snapshot_at" timestamptz NOT NULL
, "term_code" text NOT NULL
, "course_code" text NOT NULL
, "campus" text NOT NULL
, "section" int NOT NULL
, "location"    text NOT NULL
, "days"        bit(7) NOT NULL
, "start_time"  time NOT NULL
, "end_time"    time NOT NULL
, FOREIGN KEY ("snapshot_at", "term_code", "course_code", "campus", "section")
  REFERENCES "section_snapshot"
);

CREATE TABLE "section_staff"
( "snapshot_at" timestamptz NOT NULL
, "term_code" text NOT NULL
, "course_code" text NOT NULL
, "campus" text NOT NULL
, "section" int NOT NULL
, "staff_id" text NOT NULL REFERENCES "staff"
, PRIMARY KEY ("snapshot_at", "term_code", "course_code", "campus", "section", "staff_id")
, FOREIGN KEY ("snapshot_at", "term_code", "course_code", "campus", "section")
  REFERENCES "section_snapshot"
);

CREATE TABLE "section_status"
( "at" timestamptz NOT NULL
, "term_code" text NOT NULL
, "course_code" text NOT NULL
, "campus" text NOT NULL
, "section" int NOT NULL
, "status" status NOT NULL
, PRIMARY KEY ("at", "term_code", "course_code", "campus", "section")
, FOREIGN KEY ("course_code", "campus") 
  REFERENCES "course" ("code", "campus")
);

CREATE TABLE "section"
( "term_code" text NOT NULL
, "course_code" text NOT NULL
, "campus" text NOT NULL
, "section" int NOT NULL
, "snapshot_at" timestamptz NOT NULL
, "status_at" timestamptz NOT NULL
, PRIMARY KEY ("term_code", "course_code", "campus", "section")
, FOREIGN KEY ("snapshot_at", "term_code", "course_code", "campus", "section")
  REFERENCES "section_snapshot"
, FOREIGN KEY ("status_at", "term_code", "course_code", "campus", "section")
  REFERENCES "section_status"
);

CREATE OR REPLACE FUNCTION section_etc_immut()
RETURNS TRIGGER LANGUAGE plpgsql AS
$$ BEGIN RAISE EXCEPTION 'immutable'; END $$;

CREATE TRIGGER "section_snapshot_immut"
BEFORE UPDATE ON "section_snapshot"
FOR EACH ROW EXECUTE FUNCTION section_etc_immut();

CREATE TRIGGER "section_schedule_immut"
BEFORE UPDATE ON "section_schedule"
FOR EACH ROW EXECUTE FUNCTION section_etc_immut();

CREATE TRIGGER "section_staff_immut"
BEFORE UPDATE ON "section_staff"
FOR EACH ROW EXECUTE FUNCTION section_etc_immut();

CREATE TRIGGER "section_status_immut"
BEFORE UPDATE ON "section_status"
FOR EACH ROW EXECUTE FUNCTION section_etc_immut();
