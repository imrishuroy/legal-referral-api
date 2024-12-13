CREATE TABLE likes (
    like_id SERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    post_id INT,
    comment_id INT,
    type VARCHAR(7) NOT NULL CHECK (type IN ('post', 'comment')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE,
    CHECK (
        (type = 'post' AND post_id IS NOT NULL AND comment_id IS NULL) OR
        (type = 'comment' AND comment_id IS NOT NULL AND post_id IS NULL)
    ),
    UNIQUE (user_id, post_id, type),
    UNIQUE (user_id, comment_id, type)
);

-- Indexes
CREATE INDEX idx_likes_user_id ON likes (user_id);
CREATE INDEX idx_likes_post_id ON likes (post_id);
CREATE INDEX idx_likes_comment_id ON likes (comment_id);