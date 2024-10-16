CREATE TABLE devices (
    device_id VARCHAR PRIMARY KEY,
    device_token VARCHAR NOT NULL,
    user_id VARCHAR NOT NULL,
    last_used_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);