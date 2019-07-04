# aepp-contracts-library

Reverse proxy for a public facing sophia compiler that stores the source code of the contracts being compiled.

## Database

The app requires a PostgreSQL database, the following is the database schema:

```
DROP TABLE IF EXISTS "public"."contracts";

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

The contracts library can be configured using a configuration file,
by the default the app will look into `/etc/aepps/contracts_library.yaml`.

The following is an example of the configuration file:

```
# Database connection string
db_url: postgres://dbuser:dbpass@localhost/contracts_library?sslmode=disable

# Listening address
aecl_address: 0.0.0.0:1905

# List of available compilers
compilers:
- url: http://aesophia-http-v320.example.com
  version: '3.2.0'
- url: http://aesophia-http-v310.example.com
  version: '3.1.0'
  is_default: "true"
- url: http://aesophia-http-v300.example.com
  version: '3.0.0'
- url: http://aesophia-http-v210.example.com
  version: '2.1.0'

# Fine tuning for the app [OPTIONAL]
tuning:
  max_body_size: 2000000 # size of the post message in bytes
  max_open_connections: 5 # maximum numbers of open db connections
  max_idle_connections: 1 # maximum numbers of idle db connections
  version_header_name: Sophia-Compiler-Version # name of the header to use to select a compiler
```

⚠️ Note that the compilers must be available via http (not https)

⚠️ the total connection to the database opened will be `max_open_connections` + `max_idle_connections`

## Automation

To retrieve the list of the available compilers in a human readable format you can use the `Accept: application/json` that will return a json reply with the list of the available compilers.
