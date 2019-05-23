DROP TABLE IF EXISTS "public"."contracts";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Table Definition
CREATE TABLE "public"."contracts" (
    "id" varchar(200) NOT NULL,
    "source" text NOT NULL,
    "compilations" int4 NOT NULL DEFAULT 1,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp NOT NULL DEFAULT now(),
    "response_code" int4,
    "response_msg" text,
    PRIMARY KEY ("id")
);

