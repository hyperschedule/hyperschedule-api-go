CREATE TABLE "course_snapshot"
( "at" timestamptz NOT NULL DEFAULT now()
, "code"        text NOT NULL
, "campus"      text NOT NULL
, "description" text NOT NULL
, "name"        text NOT NULL
, PRIMARY KEY ("at", "code", "campus")
);

CREATE TABLE "course"
( "code" text NOT NULL
, "campus" text NOT NULL
, "snapshot_at" timestamptz NOT NULL
, PRIMARY KEY ("code", "campus")
, FOREIGN KEY ("snapshot_at", "code", "campus")
  REFERENCES "course_snapshot"
);

CREATE OR REPLACE FUNCTION course_snapshot_immut()
RETURNS TRIGGER LANGUAGE plpgsql AS 
$$ BEGIN RAISE EXCEPTION 'immutable'; END $$;

CREATE TRIGGER "course_snapshot_immut"
BEFORE UPDATE ON "course_snapshot"
FOR EACH ROW EXECUTE FUNCTION course_snapshot_immut();
