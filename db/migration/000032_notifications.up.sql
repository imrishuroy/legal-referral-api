CREATE TABLE notifications (
    notification_id SERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL, -- The user who receives the notification
    sender_id VARCHAR NOT NULL, -- The user who triggered the notification (e.g., the one who liked or commented)
    target_id INT NOT NULL, -- The ID of the entity being referenced (e.g., post_id, comment_id)
    target_type VARCHAR NOT NULL, -- Describes the type of the target (e.g., 'post', 'comment', 'share')
    notification_type VARCHAR NOT NULL, -- Type of the notification (e.g., 'like', 'comment', 'share', etc.)
    message TEXT NOT NULL, -- A description of the notification
    is_read BOOLEAN NOT NULL DEFAULT FALSE, -- Tracks if the notification has been read
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (sender_id) REFERENCES users(user_id)
);

-- Indexes to optimize queries
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);
