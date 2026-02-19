-- name: CreateUser :one

INSERT INTO users(id, name, email, password, created_at, updated_at)
VALUES($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC;