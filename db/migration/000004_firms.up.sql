CREATE TABLE firms (
    firm_id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    logo_url VARCHAR NOT NULL,
    org_type VARCHAR NOT NULL,
    website VARCHAR NOT NULL,
    location VARCHAR NOT NULL,
    about TEXT NOT NULL
);