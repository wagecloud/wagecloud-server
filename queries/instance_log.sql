-- name: GetInstanceLog :one
SELECT log.*
FROM "instance"."log" log
WHERE id = $1;

-- name: CountInstanceLogs :one
SELECT COUNT(id)
FROM "instance"."log"
WHERE (
  (instance_id = sqlc.narg('instance_id') OR sqlc.narg('instance_id') IS NULL) AND
  (type = sqlc.narg('type') OR sqlc.narg('type') IS NULL) AND
  (title ILIKE '%' || sqlc.narg('title') || '%' OR sqlc.narg('title') IS NULL) AND
  (description ILIKE '%' || sqlc.narg('description') || '%' OR sqlc.narg('description') IS NULL) AND
  (created_at >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
  (created_at <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL)
);

-- name: ListInstanceLogs :many
SELECT log.*
FROM "instance"."log" log
WHERE (
  (instance_id = sqlc.narg('instance_id') OR sqlc.narg('instance_id') IS NULL) AND
  (type = sqlc.narg('type') OR sqlc.narg('type') IS NULL) AND
  (title ILIKE '%' || sqlc.narg('title') || '%' OR sqlc.narg('title') IS NULL) AND
  (description ILIKE '%' || sqlc.narg('description') || '%' OR sqlc.narg('description') IS NULL) AND
  (created_at >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
  (created_at <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL)
)
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateInstanceLog :one
INSERT INTO "instance"."log" (instance_id, type, title, description)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateInstanceLog :one
UPDATE "instance"."log"
SET
    instance_id = COALESCE(sqlc.narg('instance_id'), instance_id),
    type = COALESCE(sqlc.narg('type'), type),
    title = COALESCE(sqlc.narg('title'), title),
    description = CASE
        WHEN sqlc.arg('null_description')::boolean THEN NULL
        ELSE COALESCE(sqlc.narg('description'), description)
    END
WHERE id = $1
RETURNING *;

-- name: DeleteInstanceLog :exec
DELETE FROM "instance"."log"
WHERE id = $1;
