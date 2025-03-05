// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: chirpy_red.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const upgradeUser = `-- name: UpgradeUser :exec
UPDATE users
SET is_chirpy_red = true, updated_at = Now()
WHERE id = $1
`

func (q *Queries) UpgradeUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, upgradeUser, id)
	return err
}
