-- name: GetNetwork :one
SELECT network.*
FROM network
WHERE id = $1;

-- name: ListNetworks :many
SELECT network.*
FROM network
WHERE (
  (id ILIKE '%' || sqlc.narg('id') || '%' OR sqlc.narg('id') IS NULL) AND
  (private_ip ILIKE '%' || sqlc.narg('private_ip') || '%' OR sqlc.narg('private_ip') IS NULL) AND
  (created_at >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
  (created_at <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL)
);

-- name: CreateNetwork :one
INSERT INTO network (id, private_ip)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateNetwork :one
UPDATE network
SET 
    id = COALESCE(sqlc.narg('new_id'), id),
    private_ip = COALESCE(sqlc.narg('private_ip'), private_ip)
WHERE id = $1
RETURNING *;

-- name: DeleteNetwork :exec
DELETE FROM network
WHERE id = $1;