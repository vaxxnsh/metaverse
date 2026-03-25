-- name: CreateElement :one
INSERT INTO elements (image_url, width, height, static)
VALUES ($1, $2, $3, $4)
RETURNING *;
