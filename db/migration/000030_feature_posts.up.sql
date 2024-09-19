CREATE TABLE feature_posts (
    feature_post_id SERIAL PRIMARY KEY,
    post_id INT NOT NULL,
    user_id VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (post_id) REFERENCES posts(post_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    UNIQUE (post_id, user_id)
);