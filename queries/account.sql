-- name: GetAccount :one
SELECT b.*, u.*
FROM account_base b
LEFT JOIN account_user u ON b.id = u.id
WHERE (
  (b.id = sqlc.narg('id') OR sqlc.narg('id') IS NULL) AND
  (b.username = sqlc.narg('username') OR sqlc.narg('username') IS NULL) AND
  (u.email = sqlc.narg('email') OR sqlc.narg('email') IS NULL)
);

-- name: CountAccounts :one
SELECT COUNT(id)
FROM account_base
WHERE (
  (id ILIKE '%' || sqlc.narg('id') || '%' OR sqlc.narg('id') IS NULL) AND
  (role = sqlc.narg('role') OR sqlc.narg('role') IS NULL) AND
  (name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL) AND
  (username ILIKE '%' || sqlc.narg('username') || '%' OR sqlc.narg('username') IS NULL) AND
  (created_at >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
  (created_at <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL)
);

-- name: ListAccounts :many
SELECT *
FROM account_base
WHERE (
  (id ILIKE '%' || sqlc.narg('id') || '%' OR sqlc.narg('id') IS NULL) AND
  (role = sqlc.narg('role') OR sqlc.narg('role') IS NULL) AND
  (name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL) AND
  (username ILIKE '%' || sqlc.narg('username') || '%' OR sqlc.narg('username') IS NULL) AND
  (created_at >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
  (created_at <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL)
)
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateAccount :one
INSERT INTO account_base (role, name, username, password)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateAccount :one
UPDATE account_base
SET
    name = COALESCE(sqlc.narg('name'), name),
    username = COALESCE(sqlc.narg('username'), username),
    password = COALESCE(sqlc.narg('password'), password)
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM account_base
WHERE id = $1;

-- name: GetUser :one
SELECT u.*, b.*
FROM account_user u
INNER JOIN account_base b ON b.id = u.id
WHERE (
  (b.role = 'USER') AND
  (b.id = sqlc.narg('id') OR sqlc.narg('id') IS NULL) AND
  (b.username = sqlc.narg('username') OR sqlc.narg('username') IS NULL) AND
  (u.email = sqlc.narg('email') OR sqlc.narg('email') IS NULL)
);

-- name: CreateUser :one
INSERT INTO account_user (id, email)
VALUES ($1, $2)
RETURNING *;