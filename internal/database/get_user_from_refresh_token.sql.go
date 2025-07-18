// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: get_user_from_refresh_token.sql

package database

import (
	"context"
)

const getUserFromRefreshToken = `-- name: GetUserFromRefreshToken :one
select id, created_at, updated_at, email, hashed_password, is_chirpy_red from users
WHERE id = (
    SELECT user_id from refresh_tokens
    WHERE token = $1
)
`

func (q *Queries) GetUserFromRefreshToken(ctx context.Context, token string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserFromRefreshToken, token)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}
