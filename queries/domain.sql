-- name: GetDomain :one
SELECT domain.*
FROM "instance"."domain" domain
WHERE id = $1;

-- name: CountDomains :one
SELECT COUNT(id)
FROM "instance"."domain"
WHERE (
  (network_id = sqlc.narg('network_id') OR sqlc.narg('network_id') IS NULL) AND
  (name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL)
);

-- name: ListDomains :many
SELECT domain.*
FROM "instance"."domain" domain
WHERE (
  (network_id = sqlc.narg('network_id') OR sqlc.narg('network_id') IS NULL) AND
  (name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL)
)
-- TODO: add order by sqlc.arg('order_by')
ORDER BY id DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateDomain :one
INSERT INTO "instance"."domain" (network_id, name)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateDomain :one
UPDATE "instance"."domain"
SET
    name = COALESCE(sqlc.narg('name'), name)
WHERE id = $1
RETURNING *;

-- name: DeleteDomain :exec
DELETE FROM "instance"."domain"
WHERE id = $1;
