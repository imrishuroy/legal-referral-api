-- name: SearchAllUsers :many
SELECT
    u.user_id,
    u.first_name,
    u.last_name,
    u.avatar_url,
    u.practice_area,
    u.practice_location
FROM
    users u
WHERE
    CONCAT(u.first_name, u.last_name) ILIKE '%' || @query::text || '%';

-- name: Search1stDegreeConnections :many
SELECT
    u.user_id,
    u.first_name,
    u.last_name,
    u.avatar_url,
    u.practice_area,
    u.practice_location
FROM
    users u
        JOIN connections c ON u.user_id = c.sender_id OR u.user_id = c.recipient_id
WHERE
    (c.sender_id = @current_user_id::text OR c.recipient_id = @current_user_id) AND
    CONCAT(u.first_name, u.last_name) ILIKE '%' || @query::text || '%'
    AND u.user_id != @current_user_id; -- Exclude the current user


-- name: Search2ndDegreeConnections :many
WITH first_degree_connections AS (
    -- Find all connections for the given user
    SELECT sender_id AS connected_user_id
    FROM connections
    WHERE recipient_id = @current_user_id::text
    UNION
    SELECT recipient_id AS connected_user_id
    FROM connections
    WHERE sender_id = @current_user_id
),
second_degree_connections AS (
    -- Find all connections of the first-degree connections
    SELECT DISTINCT c.sender_id AS second_degree_connected_user_id
    FROM connections c
    JOIN first_degree_connections fdc ON c.recipient_id = fdc.connected_user_id
    UNION
    SELECT DISTINCT c.recipient_id AS second_degree_connected_user_id
    FROM connections c
    JOIN first_degree_connections fdc ON c.sender_id = fdc.connected_user_id
)
-- Retrieve user information for the second-degree connections
SELECT u.user_id,
       u.first_name,
       u.last_name,
       u.avatar_url,
       u.practice_area,
       u.practice_location
FROM second_degree_connections sdc
JOIN users u ON sdc.second_degree_connected_user_id = u.user_id
WHERE CONCAT(u.first_name, u.last_name) ILIKE '%' || @query::text || '%'
AND u.user_id != @current_user_id; -- Exclude the current user






