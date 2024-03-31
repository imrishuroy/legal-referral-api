CREATE TABLE "license" (
    "id" bigserial PRIMARY KEY,
    "user_id" VARCHAR NOT NULL,
    "name" VARCHAR NOT NULL,
    "license_number" VARCHAR NOT NULL,
    "issue_date" VARCHAR NOT NULL,
    "issue_state" VARCHAR NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES "users"("id")
);