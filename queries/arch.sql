-- name: GetArch :one
SELECT arch.*
FROM "os"."arch" arch
WHERE id = $1;

-- name: CountArchs :one
SELECT COUNT(id)
FROM "os"."arch"
WHERE (
  (id ILIKE '%' || sqlc.narg('id') || '%' OR sqlc.narg('id') IS NULL) AND
  (name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL) AND
  (created_at >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
  (created_at <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL)
);

-- name: ListArchs :many
SELECT arch.*
FROM "os"."arch" arch
WHERE (
  (id ILIKE '%' || sqlc.narg('id') || '%' OR sqlc.narg('id') IS NULL) AND
  (name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL) AND
  (created_at >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
  (created_at <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL)
)
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateArch :one
INSERT INTO "os"."arch" (id, name)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateArch :one
UPDATE "os"."arch"
SET
    id = COALESCE(sqlc.narg('new_id'), id),
    name = COALESCE(sqlc.narg('name'), name)
WHERE id = $1
RETURNING *;

-- name: DeleteArch :exec
DELETE FROM "os"."arch"
WHERE id = $1;