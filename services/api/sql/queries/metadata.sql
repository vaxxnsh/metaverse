-- name: FindAvatarByID :one
SELECT *
FROM avatars
WHERE id = $1
LIMIT 1;