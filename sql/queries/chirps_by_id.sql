-- name: ChirpsById :one
SELECT * FROM chirps
WHERE id = $1;