CREATE TABLE socials (
    social_id  BIGSERIAL PRIMARY KEY,
    entity_id VARCHAR NOT NULL, -- This will store either user_id or company_id
    entity_type VARCHAR NOT NULL, -- 'user' or 'company'
    platform VARCHAR NOT NULL,
    link VARCHAR NOT NULL,
    UNIQUE(entity_id, entity_type, platform)
);