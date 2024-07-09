CREATE TABLE discussions (
    discussion_id SERIAL PRIMARY KEY,
    author_id     VARCHAR     NOT NULL,
    topic         TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (author_id) REFERENCES users (user_id)
);
