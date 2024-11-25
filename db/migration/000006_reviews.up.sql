CREATE TABLE reviews (
    review_id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL REFERENCES users(user_id),
    reviewer_id VARCHAR NOT NULL REFERENCES users(user_id),
    review TEXT NOT NULL,
    rating DOUBLE PRECISION NOT NULL CHECK (rating >= 1.0 AND rating <= 5.0),
    timestamp TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_user_id_reviews ON reviews(user_id);
CREATE INDEX idx_reviewer_id_reviews ON reviews(reviewer_id);
CREATE INDEX idx_timestamp_reviews ON reviews(timestamp);
