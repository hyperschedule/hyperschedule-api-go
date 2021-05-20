CREATE TABLE "course"
( "id" uuid PRIMARY KEY DEFAULT gen_random_uuid()
, "department" text NOT NULL
, "code" text NOT NULL
, "campus" text NOT NULL
, UNIQUE ("department", "code", "campus")
); -- { ("MATH", "131", "HM"): "98271ce3-4471-49ce-bd97-5cb2d8424264" }

CREATE TABLE "course_snapshot"
( "course_id" uuid NOT NULL REFERENCES "course"
, "time" timestamptz NOT NULL DEFAULT now()
, "name" text NOT NULL
, "description" text NOT NULL
, PRIMARY KEY ("course_id", "time")
);

CREATE TABLE "course_delete"
( "course_id" uuid NOT NULL
, "time" timestamptz NOT NULL DEFAULT now()
, PRIMARY KEY ("course_id", "time")
);

CREATE TABLE "course_latest"
( "course_id" uuid PRIMARY KEY REFERENCES "course"
, "snapshot_time" timestamptz NOT NULL
, FOREIGN KEY ("course_id", "snapshot_time")
  REFERENCES "course_snapshot"
);

CREATE OR REPLACE FUNCTION course_immut()
RETURNS TRIGGER LANGUAGE 'plpgsql' 
AS $$ BEGIN RAISE EXCEPTION 'course_immut'; END $$;

CREATE TRIGGER "course_immut"
BEFORE UPDATE ON "course"
EXECUTE FUNCTION course_immut();

CREATE TRIGGER "course_snapshot_immut"
BEFORE UPDATE ON "course_snapshot"
EXECUTE FUNCTION course_immut();

CREATE TRIGGER "course_delete_immut"
BEFORE UPDATE ON "course_delete"
EXECUTE FUNCTION course_immut();
