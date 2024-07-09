CREATE TABLE discussion_messages (
    message_id SERIAL PRIMARY KEY,
    parent_message_id INTEGER,
    discussion_id INTEGER NOT NULL,
    sender_id VARCHAR NOT NULL,
    message TEXT NOT NULL,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (discussion_id) REFERENCES discussions (discussion_id),
    FOREIGN KEY (sender_id) REFERENCES users (user_id)
);