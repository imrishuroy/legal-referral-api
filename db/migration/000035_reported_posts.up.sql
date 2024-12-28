CREATE TABLE reported_posts (
    report_id SERIAL PRIMARY KEY,
    post_id INT NOT NULL,
    reported_by VARCHAR(255) NOT NULL,
    reason_id INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (post_id) REFERENCES posts(post_id),
    FOREIGN KEY (reported_by) REFERENCES users(user_id),
    FOREIGN KEY (reason_id) REFERENCES report_reasons(reason_id)
);
