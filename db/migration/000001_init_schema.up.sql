CREATE TABLE "users" (
  "id" varchar PRIMARY KEY,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "mobile_number" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "bar_licence_no" varchar NOT NULL,
  "practicing_field" varchar NOT NULL,
  "experience" int NOT NULL,
  "join_date" timestamptz NOT NULL DEFAULT (now())
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

COMMENT ON COLUMN "users"."experience" IS 'in future experiences will have its own table';

ALTER TABLE "profile" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
