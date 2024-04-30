CREATE TABLE canceled_recommendations (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    recommended_user_id VARCHAR NOT NULL,
    canceled_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);
