CREATE TABLE educations (
    education_id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    school VARCHAR NOT NULL,
    degree VARCHAR NOT NULL,
    field_of_study VARCHAR NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    current BOOLEAN NOT NULL DEFAULT FALSE,
    grade VARCHAR NOT NULL,
    achievements TEXT NOT NULL,
    skills TEXT[] NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    CHECK (start_date < end_date)
);

-- Indexes
CREATE INDEX idx_user_id_educations ON educations (user_id);

