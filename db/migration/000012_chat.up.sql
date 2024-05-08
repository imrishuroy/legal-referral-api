CREATE TABLE messages (
    message_id SERIAL PRIMARY KEY,
    parent_message_id SERIAL,
    sender_id VARCHAR NOT NULL,
    recipient_id VARCHAR NOT NULL,
    message TEXT NOT NULL,
    has_attachment BOOLEAN NOT NULL DEFAULT false,
    attachment_id SERIAL,
    is_read BOOLEAN NOT NULL DEFAULT false,
    room_id VARCHAR NOT NULL,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE TABLE attachments (
    attachment_id SERIAL PRIMARY KEY,
    message_id SERIAL NOT NULL,
    attachment_url VARCHAR NOT NULL,
    attachment_type VARCHAR NOT NULL,
    FOREIGN KEY (message_id) REFERENCES messages(message_id) ON DELETE CASCADE
);

ALTER TABLE attachments
    ADD CONSTRAINT fk_attachment_message
        FOREIGN KEY (message_id) REFERENCES messages(message_id);

ALTER TABLE messages
    ADD CONSTRAINT fk_message_sender
        FOREIGN KEY (sender_id) REFERENCES users(user_id);

ALTER TABLE messages
    ADD CONSTRAINT fk_message_recipient
        FOREIGN KEY (recipient_id) REFERENCES users(user_id);

ALTER TABLE messages
    ADD CONSTRAINT unique_room_message UNIQUE (room_id, message_id);

-- Indexes
CREATE INDEX idx_sender_id ON messages(sender_id);
CREATE INDEX idx_recipient_id ON messages(recipient_id);
CREATE INDEX idx_room_id ON messages(room_id);

-- parent_message_id (handles nulls)
CREATE INDEX idx_parent_message_id ON messages(parent_message_id) WHERE parent_message_id IS NOT NULL;

-- attachment_id (handles nulls)
CREATE INDEX idx_attachment_id ON messages(attachment_id) WHERE attachment_id IS NOT NULL;

-- message_id in attachments table
CREATE INDEX idx_message_id ON attachments(message_id) WHERE message_id IS NOT NULL;