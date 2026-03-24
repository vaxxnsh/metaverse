-- name: GetAvatars :many
SELECT *
FROM avatars;

-- name: FindAvatarByID :one
SELECT *
FROM avatars
WHERE id = $1
LIMIT 1;

-- name: GetBulkUserAvatars :many
SELECT u.id AS user_id, a.image_url, a.name AS avatar_name
FROM users u
JOIN avatars a ON u.avatar_id = a.id
WHERE u.id = ANY($1::uuid[]);