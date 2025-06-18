-- name: GetAccount :one
SELECT b.*, u.*
FROM "account"."base" b
LEFT JOIN "account"."user" u ON b.id = u.id
WHERE (
  (b.type = sqlc.arg('type')) AND
  (b.id = sqlc.narg('id') OR
  b.username = sqlc.narg('username') OR
  u.email = sqlc.narg('email') OR
  u.phone = sqlc.narg('phone'))
);

-- name: CountAccounts :one
SELECT COUNT(id)
FROM "account"."base"
WHERE (
  (id ILIKE '%' || sqlc.narg('id') || '%' OR sqlc.narg('id') IS NULL) AND
  (type = sqlc.narg('type') OR sqlc.narg('type') IS NULL) AND
  (username ILIKE '%' || sqlc.narg('username') || '%' OR sqlc.narg('username') IS NULL) AND
  (created_at >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
  (created_at <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL)
);

-- name: ListAccounts :many
SELECT *
FROM "account"."base"
WHERE (
  (id ILIKE '%' || sqlc.narg('id') || '%' OR sqlc.narg('id') IS NULL) AND
  (type = sqlc.narg('type') OR sqlc.narg('type') IS NULL) AND
  (username ILIKE '%' || sqlc.narg('username') || '%' OR sqlc.narg('username') IS NULL) AND
  (created_at >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
  (created_at <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL)
)
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateAccount :one
INSERT INTO "account"."base" (type, username, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateAccount :one
UPDATE "account"."base"
SET
    username = COALESCE(sqlc.narg('username'), username),
    password = COALESCE(sqlc.narg('password'), password)
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM "account"."base"
WHERE id = $1;

-- name: GetUser :one
SELECT u.*, b.*
FROM "account"."user" u
INNER JOIN "account"."base" b ON b.id = u.id
WHERE (
  (b.type = 'ACCOUNT_TYPE_USER') AND
  (b.id = sqlc.narg('id') OR sqlc.narg('id') IS NULL) AND
  (b.username = sqlc.narg('username') OR sqlc.narg('username') IS NULL) AND
  (u.email = sqlc.narg('email') OR sqlc.narg('email') IS NULL) AND
  (u.phone = sqlc.narg('phone') OR sqlc.narg('phone') IS NULL)
);

-- name: CreateUser :one
INSERT INTO "account"."user" (id, first_name, last_name, email, phone, company, address)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateUser :one
UPDATE "account"."user"
SET
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    email = CASE
        WHEN sqlc.arg('null_email')::boolean THEN NULL
        ELSE COALESCE(sqlc.narg('email'), email)
    END,
    phone = CASE
        WHEN sqlc.arg('null_phone')::boolean THEN NULL
        ELSE COALESCE(sqlc.narg('phone'), phone)
    END,
    company = CASE
        WHEN sqlc.arg('null_company')::boolean THEN NULL
        ELSE COALESCE(sqlc.narg('company'), company)
    END,
    address = CASE
        WHEN sqlc.arg('null_address')::boolean THEN NULL
        ELSE COALESCE(sqlc.narg('address'), address)
    END
WHERE id = $1
RETURNING *;
