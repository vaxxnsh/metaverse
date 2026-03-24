-- name: GetSpaceByID :one
SELECT id, width, height FROM spaces WHERE id = $1;

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
