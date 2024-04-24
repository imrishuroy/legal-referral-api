CREATE TABLE experiences (
    experience_id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    title VARCHAR NOT NULL,
    practice_area VARCHAR NOT NULL,
    firm_id BIGSERIAL NOT NULL,
    practice_location VARCHAR NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    current BOOLEAN NOT NULL DEFAULT FALSE,
    description TEXT NOT NULL,
    skills TEXT[] NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (firm_id) REFERENCES firms(firm_id),
    CHECK (start_date < end_date)
);

-- Indexes
CREATE INDEX idx_user_id_experiences ON experiences (user_id);
CREATE INDEX idx_firm_experiences ON experiences (firm_id);
CREATE INDEX idx_practice_area_experiences ON experiences (practice_area);