-- name: CreateSpace :one
INSERT INTO spaces (name, width, height, thumbnail)
VALUES ($1, $2, $3, $4)
RETURNING *;