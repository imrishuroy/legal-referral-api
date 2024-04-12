CREATE TABLE reviews (
    review_id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL REFERENCES users(user_id), -- The user being reviewed
    reviewer_id VARCHAR NOT NULL REFERENCES users(user_id), -- The user who wrote the review
    review TEXT NOT NULL,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    timestamp TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_id_reviews ON reviews(user_id);
CREATE INDEX idx_reviewer_id_reviews ON reviews(reviewer_id);
CREATE INDEX idx_timestamp_reviews ON reviews(timestamp);