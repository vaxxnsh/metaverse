-- name: CreateSpace :one
INSERT INTO spaces (creator_id, name, width, height, thumbnail)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetMapElementsByMapID :many
SELECT * FROM map_elements WHERE map_id = $1;

-- name: CreateSpaceElement :one
INSERT INTO space_elements (space_id, element_id, x, y)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: BulkCreateSpaceElements :exec
INSERT INTO space_elements (space_id, element_id, x, y)
SELECT sqlc.arg(space_id)::uuid, unnest(sqlc.arg(element_ids)::uuid[]), unnest(sqlc.arg(xs)::int4[]), unnest(sqlc.arg(ys)::int4[]);

-- name: DeleteSpace :exec
DELETE FROM spaces WHERE id = $1;

-- name: GetSpacesByCreator :many
SELECT id, name, width, height, thumbnail FROM spaces WHERE creator_id = $1;
