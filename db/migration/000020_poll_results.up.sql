CREATE TABLE poll_results (
    poll_result_id SERIAL PRIMARY KEY,
    poll_id INT NOT NULL,
    option_index INT NOT NULL,
    user_id VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (poll_id) REFERENCES polls(poll_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    UNIQUE (poll_id, user_id)
);