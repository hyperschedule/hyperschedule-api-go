CREATE TABLE "term"
( "code"       text PRIMARY KEY -- e.g. "SU2020S1"
, "semester"   text NOT NULL    -- e.g. "SU2020"
, "date_start" date NOT NULL
, "date_end"   date NOT NULL
);
