-- name: UpdateUser :one

UPDATE users
SET email = $1, hashed_password = $2, updated_at = Now()
WHERE id = $3
RETURNING *;