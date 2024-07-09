-- ENUM type for discussion_invite_status
CREATE TYPE discussion_invite_status AS ENUM ('pending', 'accepted', 'rejected');

CREATE TABLE discussion_invites (
    discussion_invite_id SERIAL PRIMARY KEY,
    discussion_id INT NOT NULL,
    invitee_user_id VARCHAR NOT NULL,
    invited_user_id VARCHAR NOT NULL,
    status discussion_invite_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (discussion_id) REFERENCES discussion(discussion_id),
    FOREIGN KEY (invitee_user_id) REFERENCES users(user_id),
    FOREIGN KEY (invited_user_id) REFERENCES users(user_id),
    UNIQUE (discussion_id, invitee_user_id, invited_user_id)
);
