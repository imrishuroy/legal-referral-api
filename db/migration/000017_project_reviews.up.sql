CREATE TABLE project_reviews (
    review_id SERIAL PRIMARY KEY,
    project_id INT NOT NULL,
    user_id VARCHAR NOT NULL,
    review TEXT NOT NULL,
    rating DECIMAL(2, 1) NOT NULL CHECK (rating >= 0 AND rating <= 5),
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (project_id) REFERENCES projects(project_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    CONSTRAINT unique_review UNIQUE (project_id, user_id)
);

-- Optionally, add indexes if you expect frequent queries on these columns:
CREATE INDEX idx_reviews_project_id ON project_reviews(project_id);
CREATE INDEX idx_reviews_author_user_id ON project_reviews(user_id);

