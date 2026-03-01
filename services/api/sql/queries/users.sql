-- name: CreateUser :one
INSERT INTO users (id, name, email, password, created_at, updated_at)
VALUES (
    sqlc.arg(id),
    sqlc.arg(name),
    sqlc.arg(email),
    sqlc.arg(password),
    sqlc.arg(created_at),
    sqlc.arg(updated_at)
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = sqlc.arg(email)
LIMIT 1;