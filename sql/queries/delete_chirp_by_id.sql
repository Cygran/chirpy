-- name: DeleteChrip :exec
DELETE FROM chirps
WHERE id = $1;