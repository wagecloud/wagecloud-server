-- name: GetVM :one
SELECT vm.*
FROM vm
WHERE (
  id = $1 AND 
  (account_id = sqlc.narg('account_id') OR sqlc.narg('account_id') IS NULL)
);

-- name: CountVMs :one
SELECT COUNT(id)
FROM vm
WHERE (
  (account_id = sqlc.narg('account_id') OR sqlc.narg('account_id') IS NULL) AND
  (network_id = sqlc.narg('network_id') OR sqlc.narg('network_id') IS NULL) AND
  (os_id = sqlc.narg('os_id') OR sqlc.narg('os_id') IS NULL) AND
  (arch_id = sqlc.narg('arch_id') OR sqlc.narg('arch_id') IS NULL) AND
  (name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL) AND
  (cpu >= sqlc.narg('cpu_from') OR sqlc.narg('cpu_from') IS NULL) AND
  (cpu <= sqlc.narg('cpu_to') OR sqlc.narg('cpu_to') IS NULL) AND
  (ram >= sqlc.narg('ram_from') OR sqlc.narg('ram_from') IS NULL) AND
  (ram <= sqlc.narg('ram_to') OR sqlc.narg('ram_to') IS NULL) AND
  (storage >= sqlc.narg('storage_from') OR sqlc.narg('storage_from') IS NULL) AND
  (storage <= sqlc.narg('storage_to') OR sqlc.narg('storage_to') IS NULL) AND
  (created_at >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
  (created_at <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL)
);

-- name: ListVMs :many
SELECT vm.*
FROM vm
WHERE (
  (account_id = sqlc.narg('account_id') OR sqlc.narg('account_id') IS NULL) AND
  (network_id = sqlc.narg('network_id') OR sqlc.narg('network_id') IS NULL) AND
  (os_id = sqlc.narg('os_id') OR sqlc.narg('os_id') IS NULL) AND
  (arch_id = sqlc.narg('arch_id') OR sqlc.narg('arch_id') IS NULL) AND
  (name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL) AND
  (cpu >= sqlc.narg('cpu_from') OR sqlc.narg('cpu_from') IS NULL) AND
  (cpu <= sqlc.narg('cpu_to') OR sqlc.narg('cpu_to') IS NULL) AND
  (ram >= sqlc.narg('ram_from') OR sqlc.narg('ram_from') IS NULL) AND
  (ram <= sqlc.narg('ram_to') OR sqlc.narg('ram_to') IS NULL) AND
  (storage >= sqlc.narg('storage_from') OR sqlc.narg('storage_from') IS NULL) AND
  (storage <= sqlc.narg('storage_to') OR sqlc.narg('storage_to') IS NULL) AND
  (created_at >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
  (created_at <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL)
)
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateVM :one
INSERT INTO vm (id, account_id, network_id, os_id, arch_id, name, cpu, ram, storage)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateVM :one
UPDATE vm
SET 
  network_id = COALESCE(sqlc.narg('network_id'), network_id),
  os_id = COALESCE(sqlc.narg('os_id'), os_id),
  arch_id = COALESCE(sqlc.narg('arch_id'), arch_id),
  name = COALESCE(sqlc.narg('name'), name),
  cpu = COALESCE(sqlc.narg('cpu'), cpu),
  ram = COALESCE(sqlc.narg('ram'), ram),
  storage = COALESCE(sqlc.narg('storage'), storage)
WHERE (
  id = $1 AND
  (account_id = sqlc.narg('account_id') OR sqlc.narg('account_id') IS NULL)
)
RETURNING *;

-- name: DeleteVM :exec
DELETE FROM vm
WHERE (
  id = $1 AND 
  (account_id = sqlc.narg('account_id') OR sqlc.narg('account_id') IS NULL)
);