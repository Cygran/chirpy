-- name: GetUserFromRefreshToken :one
SELECT id, email, created_at, updated_at
FROM users
WHERE id = (
    SELECT user_id FROM refresh_tokens WHERE token = $1
);