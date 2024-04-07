CREATE TABLE users (
    user_id VARCHAR PRIMARY KEY NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    mobile VARCHAR,
    address VARCHAR,
    image_url VARCHAR,
    email_verified BOOLEAN NOT NULL DEFAULT false,
    mobile_verified BOOLEAN NOT NULL DEFAULT false,
    wizard_step INTEGER NOT NULL DEFAULT 0,
    wizard_completed BOOLEAN NOT NULL DEFAULT false,
    signup_method INTEGER NOT NULL,
    join_date TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE INDEX ON "users" ("first_name");

CREATE INDEX ON "users" ("last_name");
