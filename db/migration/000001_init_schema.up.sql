CREATE TABLE users (
    id VARCHAR PRIMARY KEY NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    mobile VARCHAR NOT NULL DEFAULT '',
    address VARCHAR NOT NULL DEFAULT '',
    is_email_verified BOOLEAN NOT NULL DEFAULT false,
    is_mobile_verified BOOLEAN NOT NULL DEFAULT false,
    wizard_step INTEGER NOT NULL DEFAULT 0,
    wizard_completed BOOLEAN NOT NULL DEFAULT false,
    sign_up_method INTEGER NOT NULL,
    join_date TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE INDEX ON "users" ("first_name");

CREATE INDEX ON "users" ("last_name");
