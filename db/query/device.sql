-- name: SaveDevice :exec
INSERT INTO devices (
    device_id,
    device_token,
    user_id,
    last_used_at
) VALUES (
    $1, $2, $3, current_timestamp
)
ON CONFLICT (device_id)
    DO UPDATE SET
    device_token = EXCLUDED.device_token,
    last_used_at = current_timestamp;
