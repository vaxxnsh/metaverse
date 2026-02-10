-- name: CreateUser :one

INSERT INTO users(id, name, password, created_at, updated_at)
VALUES($1,$2,$3,$4,$5)
RETURNING id, name, created_at, updated_at;

-- name: ListUsers :many
SELECT id, name, created_at, updated_at
FROM users
ORDER BY created_at DESC;