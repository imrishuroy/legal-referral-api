CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id VARCHAR(20) PRIMARY KEY DEFAULT ('USR' || substring(encode(gen_random_bytes(8), 'hex') from 1 for 8)),
    email VARCHAR UNIQUE NOT NULL,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    is_email_verified BOOLEAN NOT NULL DEFAULT false,
    join_date TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);


CREATE TABLE "profile" (
  "id" bigserial PRIMARY KEY,
  "user_id" varchar NOT NULL,
  "headline" varchar NOT NULL,
  "summary" varchar NOT NULL,
  "industry" varchar NOT NULL,
  "website" varchar NOT NULL
);

CREATE INDEX ON "users" ("first_name");

CREATE INDEX ON "users" ("last_name");

-- COMMENT ON COLUMN "users"."experience" IS 'in future experiences will have its own table';

ALTER TABLE "profile" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
