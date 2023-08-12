-- name: RegisterUser :one
INSERT INTO user_credentials (id, password_hash, email)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM user_credentials WHERE email = $1;