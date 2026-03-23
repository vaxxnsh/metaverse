-- name: CreateUser :one
INSERT INTO users (name, email, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: FindUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: PatchUserMetadata :one
UPDATE users
SET avatar_id = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
