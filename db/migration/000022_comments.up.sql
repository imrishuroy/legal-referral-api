CREATE TABLE comments (
    comment_id SERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    post_id INT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    parent_comment_id INT,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
    FOREIGN KEY (parent_comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_post_id ON comments (post_id);
CREATE INDEX idx_parent_comment_id ON comments (parent_comment_id);
