-- ENUM type for post_type
CREATE TYPE post_type AS ENUM ('text', 'image', 'video', 'audio', 'link', 'document', 'poll', 'other');

CREATE TABLE posts (
    post_id SERIAL PRIMARY KEY,
    owner_id VARCHAR NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    media TEXT[] NOT NULL DEFAULT '{}',
    post_type post_type NOT NULL DEFAULT 'text',
    poll_id INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);