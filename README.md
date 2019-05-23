# aepp-contracts-library

Reverse proxy for a public facing sophia compiler that stores the source code of the contracts beein compiled

## Schema

The following is the database schema:

```
DROP TABLE IF EXISTS "public"."contracts";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Table Definition
CREATE TABLE "public"."contracts" (
    "id" varchar(200) NOT NULL, -- blacke2b hash of the contract source code
    "source" text NOT NULL, -- the source code of the contract
    "compilations" int4 NOT NULL DEFAULT 1, -- number of compilations
    "created_at" timestamp NOT NULL DEFAULT now(), -- first compilation date
    "updated_at" timestamp NOT NULL DEFAULT now(), -- the last time the contract was compiled
    "response_code" int4, -- the response code for the compilation (for the first call)
    "response_msg" text,  -- the response message for the compilation (for the first call)
    PRIMARY KEY ("id")
);
```

## Configuration

The contracts library can be configured using environment variables:

#### `COMPILER_URL`

The url of the compiler, for example: `http://compiler.aepps.com`.

Note that the compiler must be available via http (not https)

#### `DATABASE_URL`

The connection string for the database, for example: `postgres://aecl:aecl@localhost/contracts_library?sslmode=disable`.

The database must be a PostgreSQL database.

#### `LISTEN_ADDRES`

The address the application listens to, default `0.0.0.0:1905`

#### `MAX_BODY_SIZE`

The maximum size of the body for the request and the response in bytes, default `2e6` (2mb)

## Building

to 

## Kubernetes 

The following is a sample kubernetes configuration


