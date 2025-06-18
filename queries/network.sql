-- name: GetNetwork :one
SELECT network.*
FROM "instance"."network" network
WHERE (
  id = sqlc.narg('id') OR instance_id = sqlc.narg('instance_id')
);

-- name: CountNetworks :one
SELECT COUNT(id)
FROM "instance"."network"
WHERE (
  (instance_id = sqlc.narg('instance_id') OR sqlc.narg('instance_id') IS NULL) AND
  (private_ip ILIKE '%' || sqlc.narg('private_ip') || '%' OR sqlc.narg('private_ip') IS NULL) AND
  (mac_address ILIKE '%' || sqlc.narg('mac_address') || '%' OR sqlc.narg('mac_address') IS NULL) AND
  (public_ip ILIKE '%' || sqlc.narg('public_ip') || '%' OR sqlc.narg('public_ip') IS NULL)
);

-- name: ListNetworks :many
SELECT network.*
FROM "instance"."network" network
WHERE (
  (instance_id = sqlc.narg('instance_id') OR sqlc.narg('instance_id') IS NULL) AND
  (private_ip ILIKE '%' || sqlc.narg('private_ip') || '%' OR sqlc.narg('private_ip') IS NULL) AND
  (mac_address ILIKE '%' || sqlc.narg('mac_address') || '%' OR sqlc.narg('mac_address') IS NULL) AND
  (public_ip ILIKE '%' || sqlc.narg('public_ip') || '%' OR sqlc.narg('public_ip') IS NULL)
)
ORDER BY id DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateNetwork :one
INSERT INTO "instance"."network" (instance_id, private_ip, mac_address, public_ip)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateNetwork :one
UPDATE "instance"."network"
SET
    private_ip = COALESCE(sqlc.narg('private_ip'), private_ip),
    mac_address = COALESCE(sqlc.narg('mac_address'), mac_address),
    public_ip = CASE
        WHEN sqlc.arg('null_public_ip')::boolean THEN NULL
        ELSE COALESCE(sqlc.narg('public_ip'), public_ip)
    END
WHERE (
  (id = sqlc.narg('id') OR sqlc.narg('id') IS NULL) AND
  (instance_id = sqlc.narg('instance_id') OR sqlc.narg('instance_id') IS NULL)
)
RETURNING *;

-- name: DeleteNetwork :exec
DELETE FROM "instance"."network"
WHERE id = $1;
