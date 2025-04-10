-- name: GetOS :one
SELECT os.*
FROM os
WHERE id = $1;

-- name: CountOSs :one
SELECT COUNT(id)
FROM os
WHERE (
  (id ILIKE '%' || sqlc.narg('id') || '%' OR sqlc.narg('id') IS NULL) AND
  (name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL) AND
  (created_at >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
  (created_at <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL)
);

-- name: ListOSs :many
SELECT os.*
FROM os
WHERE (
  (id ILIKE '%' || sqlc.narg('id') || '%' OR sqlc.narg('id') IS NULL) AND
  (name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL) AND
  (created_at >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
  (created_at <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL)
)
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateOS :one
INSERT INTO os (name)
VALUES ($1)
RETURNING *;

-- name: UpdateOS :one
UPDATE os
SET
    id = COALESCE(sqlc.narg('new_id'), id),
    name = COALESCE(sqlc.narg('name'), name)
WHERE id = $1
RETURNING *;

-- name: DeleteOS :exec
DELETE FROM os
WHERE id = $1;
