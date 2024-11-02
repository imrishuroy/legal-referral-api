-- name: CommentPost :one
INSERT INTO comments (
    user_id,
    post_id,
    content,
    parent_comment_id
) VALUES (
    $1, $2, $3, $4
) RETURNING
    comment_id,
    user_id as author_user_id,
    post_id,
    content,
    created_at,
    parent_comment_id;

-- name: ListComments :many
WITH RECURSIVE comment_tree AS (
    SELECT
        c.comment_id,
        c.user_id,
        c.post_id,
        c.content,
        c.created_at,
        c.parent_comment_id
    FROM
        comments c
    WHERE
        c.post_id = $1 AND c.parent_comment_id IS NULL
    UNION ALL
    SELECT
        c.comment_id,
        c.user_id,
        c.post_id,
        c.content,
        c.created_at,
        c.parent_comment_id
    FROM
        comments c
            INNER JOIN
        comment_tree ct ON ct.comment_id = c.parent_comment_id
)
SELECT
    comment_tree.comment_id,
    comment_tree.post_id,
    comment_tree.content,
    comment_tree.created_at,
    comment_tree.parent_comment_id,
    comment_tree.user_id AS author_user_id,
    u.first_name AS author_first_name,
    u.last_name AS author_last_name,
    u.avatar_url AS author_avatar_url,
    u.practice_area AS author_practice_area
FROM
    comment_tree
        JOIN
    users u ON comment_tree.user_id = u.user_id;

-- name: ListComments2 :many
WITH RECURSIVE comment_tree AS (
    SELECT
        c.comment_id,
        c.user_id,
        c.post_id,
        c.content,
        c.created_at,
        c.parent_comment_id
    FROM
        comments c
    WHERE
        c.post_id = $1 AND c.parent_comment_id IS NULL
    UNION ALL
    SELECT
        c.comment_id,
        c.user_id,
        c.post_id,
        c.content,
        c.created_at,
        c.parent_comment_id
    FROM
        comments c
            INNER JOIN
        comment_tree ct ON ct.comment_id = c.parent_comment_id
)
SELECT
    comment_tree.comment_id,
    comment_tree.post_id,
    comment_tree.content,
    comment_tree.created_at,
    comment_tree.parent_comment_id,
    comment_tree.user_id AS author_user_id,
    u.first_name AS author_first_name,
    u.last_name AS author_last_name,
    u.avatar_url AS author_avatar_url,
    u.practice_area AS author_practice_area,
    COALESCE(likes_count_table.likes_count, 0) AS likes_count,
    CASE
        WHEN liked_comments.comment_id IS NOT NULL THEN true
        ELSE false
        END AS is_liked
FROM
    comment_tree
        JOIN
    users u ON comment_tree.user_id = u.user_id
        LEFT JOIN
    (
        SELECT
            l.comment_id,
            COUNT(*) AS likes_count
        FROM
            likes l
        WHERE
            l.type = 'comment'
        GROUP BY
            l.comment_id
    ) AS likes_count_table ON comment_tree.comment_id = likes_count_table.comment_id
        LEFT JOIN
    (
        SELECT
            l.comment_id
        FROM
            likes l
        WHERE
            l.type = 'comment' AND l.user_id = $2
    ) AS liked_comments ON comment_tree.comment_id = liked_comments.comment_id
ORDER BY
    comment_tree.created_at;

-- name: GetPostCommentsCount :one
SELECT
    COUNT(*) AS comments_count
FROM comments
WHERE post_id = $1;