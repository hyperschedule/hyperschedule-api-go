CREATE TABLE "staff"
( "id"         text PRIMARY KEY DEFAULT gen_random_uuid()
, "lingk_id"   text NOT NULL UNIQUE
, "first_name" text NOT NULL
, "last_name"  text NOT NULL
);
