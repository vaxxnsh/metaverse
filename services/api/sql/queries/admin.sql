-- name: CreateAdmin :one
INSERT INTO admins (name, email, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: FindAdminByEmail :one
SELECT *
FROM admins
WHERE email = $1
LIMIT 1;

-- name: FindAdminByID :one
SELECT *
FROM admins
WHERE id = $1
LIMIT 1;

-- name: PatchAdminMetadata :one
UPDATE admins
SET avatar_id = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;