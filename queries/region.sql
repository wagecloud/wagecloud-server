-- name: GetRegion :one
SELECT region.*
FROM "instance"."region" region
WHERE id = $1;

-- name: CountRegions :one
SELECT COUNT(id)
FROM "instance"."region"
WHERE (
  (id = sqlc.narg('id') OR sqlc.narg('id') IS NULL) AND
  (name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL)
);

-- name: ListRegions :many
SELECT region.*
FROM "instance"."region" region
WHERE (
  (id = sqlc.narg('id') OR sqlc.narg('id') IS NULL) AND
  (name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL)
)
ORDER BY id DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateRegion :one
INSERT INTO "instance"."region" (id, name)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateRegion :one
UPDATE "instance"."region"
SET
    id = COALESCE(sqlc.narg('new_id'), id),
    name = COALESCE(sqlc.narg('name'), name)
WHERE id = $1
RETURNING *;

-- name: DeleteRegion :exec
DELETE FROM "instance"."region"
WHERE id = $1;
