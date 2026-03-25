-- name: GetAllElements :many
SELECT id, image_url, width, height, static FROM elements;

-- name: GetSpaceByID :one
SELECT id, width, height FROM spaces WHERE id = $1;

-- name: DeleteSpaceElement :exec
DELETE FROM space_elements WHERE space_id = $1 AND x = $2 AND y = $3;

-- name: GetSpaceElements :many
SELECT
    se.x,
    se.y,
    e.id       AS element_id,
    e.image_url,
    e.width    AS element_width,
    e.height   AS element_height,
    e.static
FROM space_elements se
JOIN elements e ON se.element_id = e.id
WHERE se.space_id = $1;
