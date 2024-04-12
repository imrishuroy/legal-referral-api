CREATE TABLE educations (
    education_id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    school VARCHAR NOT NULL,
    degree VARCHAR NOT NULL,
    field_of_study VARCHAR NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ,
    current BOOLEAN NOT NULL DEFAULT FALSE,
    grade VARCHAR NOT NULL,
    achievements TEXT,
    skills TEXT[] NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    CHECK (start_date < end_date),
    CHECK (current = (end_date IS NULL))
);

-- Indexes
CREATE INDEX idx_user_id_educations ON educations (user_id);

