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

-- name: CreateMap :one
INSERT INTO maps (name, width, height)
VALUES ($1, $2, $3)
RETURNING *;

-- name: BulkCreateMapElements :exec
INSERT INTO map_elements (map_id, element_id, x, y)
SELECT sqlc.arg(map_id)::uuid, unnest(sqlc.arg(element_ids)::uuid[]), unnest(sqlc.arg(xs)::int4[]), unnest(sqlc.arg(ys)::int4[]);
