-- name: CreateUser :one
INSERT INTO users (id, display_name)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;