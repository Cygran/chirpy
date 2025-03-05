-- name: UpgradeUser :exec
UPDATE users
SET is_chirpy_red = true, updated_at = Now()
WHERE id = $1;