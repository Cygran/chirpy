-- name: GetRefreshToken :one
SELECT token, expires_at, revoked_at FROM refresh_tokens WHERE token = $1;