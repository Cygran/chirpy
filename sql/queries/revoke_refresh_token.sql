-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = Now(), updated_at = NOW()
WHERE token = $1;