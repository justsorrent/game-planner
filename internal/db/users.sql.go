// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: users.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, display_name)
VALUES ($1, $2)
RETURNING id, display_name, created_at, updated_at
`

type CreateUserParams struct {
	ID          uuid.UUID
	DisplayName string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.ID, arg.DisplayName)
	var i User
	err := row.Scan(
		&i.ID,
		&i.DisplayName,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, display_name, created_at, updated_at FROM users
WHERE id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.DisplayName,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
