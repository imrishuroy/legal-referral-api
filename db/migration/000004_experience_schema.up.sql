CREATE TABLE experiences (
    "id" bigserial PRIMARY KEY,
    "user_id" VARCHAR UNIQUE NOT NULL,
    "practice_area" VARCHAR NOT NULL,
    "practice_location" VARCHAR NOT NULL,
    "experience" int NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);