-- name: GetNetwork :one
SELECT network.*
FROM "instance"."network" network
WHERE id = $1;

-- name: CountNetworks :one
SELECT COUNT(id)
FROM "instance"."network"
WHERE (
  (id ILIKE '%' || sqlc.narg('id') || '%' OR sqlc.narg('id') IS NULL) AND
  (private_ip ILIKE '%' || sqlc.narg('private_ip') || '%' OR sqlc.narg('private_ip') IS NULL) AND
  (created_at >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
  (created_at <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL)
);

-- name: ListNetworks :many
SELECT network.*
FROM "instance"."network" network
WHERE (
  (id ILIKE '%' || sqlc.narg('id') || '%' OR sqlc.narg('id') IS NULL) AND
  (private_ip ILIKE '%' || sqlc.narg('private_ip') || '%' OR sqlc.narg('private_ip') IS NULL) AND
  (created_at >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
  (created_at <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL)
)
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateNetwork :one
INSERT INTO "instance"."network" (id, private_ip)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateNetwork :one
UPDATE "instance"."network"
SET
    private_ip = COALESCE(sqlc.narg('private_ip'), private_ip)
WHERE id = $1
RETURNING *;

-- name: DeleteNetwork :exec
DELETE FROM "instance"."network"
WHERE id = $1;
