

-- name: CreateAccount :one
INSERT INTO account (id, name, email, created_at)
VALUES (DEFAULT, $1, $2, NOW())
RETURNING id, name, email, created_at;


-- name: GetAccountByID :one
SELECT id, name, email, created_at
FROM account
WHERE id = $1;


-- name: UpdateAccount :one
UPDATE account
SET name = $2, email = $3
WHERE id = $1
RETURNING id, name, email, created_at;


-- name: DeleteAccount :exec
DELETE FROM account
WHERE id = $1;

-- name: ListAccounts :many
SELECT id, name, email, created_at
FROM account
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
