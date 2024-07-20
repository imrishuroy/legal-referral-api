-- name: CreateDiscussion :one
INSERT INTO discussions (
    author_id,
    topic
) VALUES (
    $1, $2
) RETURNING *;

-- name: UpdateDiscussionTopic :exec
UPDATE discussions
SET topic = $2
WHERE discussion_id = $1;

-- name: InviteUserToDiscussion :exec
INSERT INTO discussion_invites (
    discussion_id,
    invitee_user_id,
    invited_user_id
) VALUES (
    $1, $2, $3
);

-- name: JoinDiscussion :exec
UPDATE discussion_invites SET status = 'accepted'
WHERE discussion_id = $1 AND invited_user_id = $2;

-- name: RejectDiscussion :exec
UPDATE discussion_invites SET status = 'rejected'
WHERE discussion_id = $1 AND invited_user_id = $2;

-- name: ListActiveDiscussions :many
SELECT
    d.discussion_id,
    d.author_id,
    d.topic,
    d.created_at,
    COUNT(DISTINCT di.invited_user_id) AS active_member_count
FROM
    discussions d
        LEFT JOIN discussion_invites di ON d.discussion_id = di.discussion_id
WHERE
    d.author_id = $1
   OR (di.invited_user_id = $1 AND di.status = 'accepted')
   OR (di.invitee_user_id = $1 AND di.status = 'accepted')
GROUP BY
    d.discussion_id,
    d.author_id,
    d.topic,
    d.created_at
ORDER BY d.created_at DESC;

-- name: ListDiscussionInvites :many
SELECT sqlc.embed(discussion_invites), sqlc.embed(discussions), sqlc.embed(users)
FROM discussion_invites
JOIN discussions ON discussion_invites.discussion_id = discussions.discussion_id
JOIN users ON discussion_invites.invitee_user_id = users.user_id
WHERE discussion_invites.invited_user_id = $1
    AND discussion_invites.status = 'pending';


-- name: ListDiscussionParticipants :many
SELECT
    u.user_id,
    u.first_name,
    u.last_name,
    u.avatar_url,
    u.practice_area
FROM
    discussion_invites di
        JOIN
    users u
    ON
        di.invited_user_id = u.user_id
WHERE
    di.status = 'accepted'
  AND di.discussion_id = $1
ORDER BY
    di.created_at;

-- name: ListUninvitedParticipants :many
WITH invited_users AS (
    SELECT di.invitee_user_id AS user_id
    FROM discussion_invites di
    WHERE di.discussion_id = $1

    UNION

    SELECT di.invited_user_id AS user_id
    FROM discussion_invites di
    WHERE di.discussion_id = $1
)
SELECT
    u.user_id,
    u.first_name,
    u.last_name,
    u.avatar_url,
    u.practice_area
FROM users u
WHERE u.user_id NOT IN (
    SELECT iu.user_id
    FROM invited_users iu
);
