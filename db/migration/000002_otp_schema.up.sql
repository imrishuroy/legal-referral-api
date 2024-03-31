CREATE TABLE otps (
  "session_id" bigserial PRIMARY KEY,
  "email" varchar NOT NULL,
  "channel" varchar NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
  "otp" int NOT NULL
);