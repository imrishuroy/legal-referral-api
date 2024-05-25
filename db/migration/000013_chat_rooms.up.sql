CREATE TABLE chat_rooms (
    room_id VARCHAR PRIMARY KEY,
    user1_id VARCHAR NOT NULL,
    user2_id VARCHAR NOT NULL,
    last_message_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    UNIQUE (user1_id, user2_id),
    FOREIGN KEY (user1_id) REFERENCES users(user_id),
    FOREIGN KEY (user2_id) REFERENCES users(user_id)

);

-- Indexs
CREATE INDEX chat_rooms_sender_id_idx ON chat_rooms(user1_id);
CREATE INDEX chat_rooms_recepient_id_idx ON chat_rooms(user2_id);
CREATE INDEX chat_rooms_last_message_at_idx ON chat_rooms(last_message_at);
CREATE INDEX chat_rooms_created_at_idx ON chat_rooms(created_at);