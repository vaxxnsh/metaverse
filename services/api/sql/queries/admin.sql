-- name: CreateAdmin :one
INSERT INTO admins (name, email, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: FindAdminByEmail :one
SELECT *
FROM admins
WHERE email = $1
LIMIT 1;