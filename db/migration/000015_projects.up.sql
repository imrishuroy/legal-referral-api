-- ENUM type for status
CREATE TYPE  project_status AS ENUM ('active', 'awarded', 'accepted', 'rejected', 'started', 'complete_initiated', 'completed', 'cancelled');

CREATE TABLE projects (
    project_id SERIAL PRIMARY KEY,
    referred_user_id VARCHAR NOT NULL,
    referrer_user_id VARCHAR NOT NULL,
    referral_id INT NOT NULL,
    status project_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    FOREIGN KEY (referred_user_id) REFERENCES users(user_id),
    FOREIGN KEY (referrer_user_id) REFERENCES users(user_id),
    CHECK (completed_at IS NULL OR started_at IS NOT NULL),
    CHECK (completed_at IS NULL OR completed_at > started_at)
);

-- Indexes for foreign keys
CREATE INDEX idx_referred_user_id ON projects(referred_user_id);
CREATE INDEX idx_referrer_user_id ON projects(referrer_user_id);