CREATE TABLE "term"
( "code"       text PRIMARY KEY -- e.g. "SU2020S1"
, "semester"   text NOT NULL    -- e.g. "SU2020"
, "start_date" date NOT NULL
, "end_date"   date NOT NULL
);
