CREATE TABLE reviews (
    review_id SERIAL PRIMARY KEY,
    user_id VARCHAR REFERENCES users(user_id), -- The user being reviewed
    reviewer_id VARCHAR REFERENCES users(user_id), -- The user who wrote the review
    review_text TEXT,
    rating INT,
    timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_id_reviews ON reviews(user_id);
CREATE INDEX idx_reviewer_id_reviews ON reviews(reviewer_id);
CREATE INDEX idx_timestamp_reviews ON reviews(timestamp);