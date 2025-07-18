// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: delete_chirp_by_id.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const deleteChirpFromId = `-- name: DeleteChirpFromId :exec
DELETE FROM chirps
WHERE id = $1
`

func (q *Queries) DeleteChirpFromId(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteChirpFromId, id)
	return err
}
