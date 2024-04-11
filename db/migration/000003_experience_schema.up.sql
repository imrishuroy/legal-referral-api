CREATE TABLE experiences (
    experience_id bigserial PRIMARY KEY,
    user_id VARCHAR UNIQUE NOT NULL,
    practice_area VARCHAR NOT NULL,
    practice_location VARCHAR NOT NULL,
    experience VARCHAR NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);