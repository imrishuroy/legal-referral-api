CREATE TABLE connection_invitations (
    id SERIAL PRIMARY KEY,
    sender_id VARCHAR NOT NULL,
    recipient_id VARCHAR NOT NULL,
    status INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    UNIQUE (sender_id, recipient_id)
);
