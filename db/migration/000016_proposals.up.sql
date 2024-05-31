-- ENUM type for status
CREATE TYPE proposal_status AS ENUM ('active', 'accepted', 'rejected', 'cancelled');

CREATE TABLE proposals (
    proposal_id SERIAL PRIMARY KEY,
    project_id INT NOT NULL,
    user_id VARCHAR NOT NULL,
    title TEXT NOT NULL,
    proposal TEXT NOT NULL,
    status proposal_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (project_id) REFERENCES projects(project_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    CONSTRAINT unique_proposal UNIQUE (project_id, user_id)
);