CREATE TABLE post_statistics (
    post_id INT PRIMARY KEY REFERENCES posts(post_id),
    views BIGINT NOT NULL DEFAULT 0,
    likes BIGINT NOT NULL DEFAULT 0,
    comments BIGINT NOT NULL DEFAULT 0,
    shares BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts (post_id) ON DELETE CASCADE
);

-- Adding individual indexes on likes and comments columns
CREATE INDEX idx_post_statistics_likes ON post_statistics (likes);
CREATE INDEX idx_post_statistics_comments ON post_statistics (comments );