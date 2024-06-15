CREATE TABLE polls (
    poll_id SERIAL PRIMARY KEY,
    owner_id VARCHAR NOT NULL,
    title TEXT NOT NULL,
    options TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    end_time TIMESTAMPTZ
);