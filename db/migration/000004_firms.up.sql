CREATE TABLE firms (
    firm_id BIGSERIAL PRIMARY KEY,
    owner_user_id VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    logo_url VARCHAR NOT NULL,
    org_type VARCHAR NOT NULL,
    website VARCHAR NOT NULL,
    location VARCHAR NOT NULL,
    about TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (owner_user_id) REFERENCES users (user_id)
);
