CREATE TABLE users (
    user_id VARCHAR PRIMARY KEY,
    email VARCHAR UNIQUE NOT NULL,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    about VARCHAR,
    mobile VARCHAR,
    address VARCHAR,
    avatar_url VARCHAR,
    banner_url VARCHAR,
    email_verified BOOLEAN NOT NULL DEFAULT false,
    mobile_verified BOOLEAN NOT NULL DEFAULT false,
    wizard_step INTEGER NOT NULL DEFAULT 0,
    wizard_completed BOOLEAN NOT NULL DEFAULT false,
    signup_method INTEGER NOT NULL,
    practice_area VARCHAR,
    practice_location VARCHAR,
    experience VARCHAR,
    average_billing_per_client INTEGER,
    case_resolution_rate INTEGER,
    open_to_referral BOOLEAN NOT NULL DEFAULT false,
    license_verified BOOLEAN NOT NULL DEFAULT false,
    license_rejected BOOLEAN NOT NULL DEFAULT false,
    join_date TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

-- Indexes
CREATE INDEX ON "users" ("first_name");
CREATE INDEX ON "users" ("last_name");
