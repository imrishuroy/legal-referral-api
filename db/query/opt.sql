-- name: StoreOTP :one
INSERT INTO otps (
    email,
    channel,
    otp
) VALUES (
   $1, $2, $3
) RETURNING session_id;

-- name: GetOTP :one
SELECT * FROM otps
WHERE session_id = $1;

-- name: DeleteOTP :exec
DELETE FROM otps
WHERE session_id = $1;
