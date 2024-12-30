-- name: ListUserPosts :many
SELECT posts.*,
       post_owner.first_name AS owner_first_name,
       post_owner.last_name AS owner_last_name,
       post_owner.avatar_url AS owner_avatar_url,
       post_owner.practice_area AS owner_practice_area,
       COALESCE(post_stats.likes, 0) AS likes_count,
       COALESCE(post_stats.comments, 0) AS comments_count
FROM posts
         JOIN users post_owner ON posts.owner_id = post_owner.user_id
         LEFT JOIN post_statistics post_stats ON posts.post_id = post_stats.post_id
WHERE posts.owner_id = $1
ORDER BY posts.created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListUserComments :many
SELECT
       comments.comment_id,
       comments.post_id,
       comments.content,
       comments.created_at,
       comments.parent_comment_id,
       users.user_id AS author_user_id,
       users.first_name AS author_first_name,
       users.last_name AS author_last_name,
       users.avatar_url AS author_avatar_url
FROM comments
         JOIN users ON comments.user_id = users.user_id
WHERE comments.user_id = $1
ORDER BY comments.created_at DESC
LIMIT $2
OFFSET $3;

-- name: GetUserFollowersCount :one
SELECT
COALESCE((SELECT COUNT(*)
          FROM connection_invitations
          WHERE recipient_id = users.user_id
            AND status NOT IN ('rejected', 'cancelled')), 0) AS followers_count
FROM users
WHERE users.user_id = $1;

