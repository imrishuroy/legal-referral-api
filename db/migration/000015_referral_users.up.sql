CREATE TABLE referral_users (
    referral_user_id SERIAL PRIMARY KEY,
    project_id INT NOT NULL,
    referred_user_id VARCHAR NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(project_id),
    FOREIGN KEY (referred_user_id) REFERENCES users(user_id),
    UNIQUE (project_id, referred_user_id)
);
