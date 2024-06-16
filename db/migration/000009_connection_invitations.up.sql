-- Define the ENUM type
CREATE TYPE invitation_status AS ENUM ('pending', 'accepted', 'rejected', 'cancelled', 'none');

CREATE TABLE connection_invitations (
    id SERIAL PRIMARY KEY,
    sender_id VARCHAR NOT NULL,
    recipient_id VARCHAR NOT NULL,
    status invitation_status NOT NULL DEFAULT 'none',
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    UNIQUE (sender_id, recipient_id)
);
