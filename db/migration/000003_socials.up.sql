CREATE TABLE socials (
    social_id  BIGSERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    platform_name VARCHAR NOT NULL,
    link_url VARCHAR NOT NULL,
    UNIQUE(user_id, platform_name)
);
