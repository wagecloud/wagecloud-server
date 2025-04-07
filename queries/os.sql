-- name: GetOS :one
SELECT os.*
FROM os
WHERE id = $1;

-- name: ListOSs :many
SELECT os.*
FROM os
WHERE (
  (name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL) AND
  (created_at >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
  (created_at <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL)
);

-- name: CreateOS :one
INSERT INTO os (name)
VALUES ($1)
RETURNING *;

-- name: UpdateOS :one
UPDATE os
SET
    name = COALESCE(sqlc.narg('name'), name)
WHERE id = $1
RETURNING *;

-- name: DeleteOS :exec
DELETE FROM os
WHERE id = $1;