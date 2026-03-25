-- name: CreateElement :one
INSERT INTO elements (image_url, width, height, static)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateElementImage :one
UPDATE elements SET image_url = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: CreateAvatar :one
INSERT INTO avatars (image_url, name)
VALUES ($1, $2)
RETURNING *;
