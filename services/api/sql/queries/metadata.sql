-- name: GetAvatars :many
SELECT *
FROM avatars;

-- name: FindAvatarByID :one
SELECT *
FROM avatars
WHERE id = $1
LIMIT 1;