-- name: ListAttorneys :many
WITH conn AS (
    SELECT
        user_id,
        COUNT(*) AS total_connections
    FROM (
             SELECT sender_id AS user_id FROM connections
             UNION ALL
             SELECT recipient_id AS user_id FROM connections
         ) AS all_connections
    GROUP BY user_id
)
SELECT
    u.user_id,
    u.first_name,
    u.last_name,
    u.practice_area,
    u.practice_location,
    p.price_id,
    p.service_type,
    p.per_hour_price,
    p.per_hearing_price,
    p.contingency_price,
    p.hybrid_price,
    COALESCE(conn.total_connections, 0) AS total_connections
FROM
    users u
        LEFT JOIN conn ON u.user_id = conn.user_id
        LEFT JOIN pricing p ON u.user_id = p.user_id
ORDER BY
    u.user_id
LIMIT $1
OFFSET $2;
