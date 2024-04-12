CREATE TABLE licenses (
    license_id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    license_number VARCHAR NOT NULL,
    issue_date TIMESTAMPTZ NOT NULL,
    issue_state VARCHAR NOT NULL,
    license_pdf VARCHAR,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- Index on user_id for foreign key constraint
CREATE INDEX ON licenses (user_id);
