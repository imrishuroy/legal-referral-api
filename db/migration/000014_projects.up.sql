-- ENUM type for status
CREATE TYPE  project_status AS ENUM ('active', 'awarded', 'accepted', 'rejected', 'started', 'complete_initiated', 'completed', 'cancelled');

CREATE TABLE projects (
    project_id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    preferred_practice_area TEXT NOT NULL,
    preferred_practice_location TEXT NOT NULL,
    case_description TEXT NOT NULL,
    referrer_user_id VARCHAR NOT NULL,
    referred_user_id VARCHAR,
    status project_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    FOREIGN KEY (referrer_user_id) REFERENCES users(user_id),
    CHECK (completed_at IS NULL OR started_at IS NOT NULL),
    CHECK (completed_at IS NULL OR completed_at > started_at),
    CONSTRAINT referrer_user_id_referral_id_unique UNIQUE (referrer_user_id, referral_id)
);

-- Optionally, add indexes if you expect frequent queries on these columns:
CREATE INDEX idx_projects_referrer_user_id ON projects(referrer_user_id);

