-- name: GetPayment :one
SELECT p.*
FROM "payment"."base" p
WHERE p.id = $1;

-- name: CountPayments :one
SELECT COUNT(id)
FROM "payment"."base"
WHERE (
  (id ILIKE '%' || sqlc.narg('id') || '%' OR sqlc.narg('id') IS NULL) AND
  (account_id = sqlc.narg('account_id') OR sqlc.narg('account_id') IS NULL) AND
  (method = sqlc.narg('method') OR sqlc.narg('method') IS NULL) AND
  (status = sqlc.narg('status') OR sqlc.narg('status') IS NULL) AND
  (date_created >= sqlc.narg('date_created_from') OR sqlc.narg('date_created_from') IS NULL) AND
  (date_created <= sqlc.narg('date_created_to') OR sqlc.narg('date_created_to') IS NULL)
);

-- name: ListPayments :many
SELECT p.*
FROM "payment"."base" p
WHERE (
  (p.id ILIKE '%' || sqlc.narg('id') || '%' OR sqlc.narg('id') IS NULL) AND
  (p.account_id = sqlc.narg('account_id') OR sqlc.narg('account_id') IS NULL) AND
  (p.method = sqlc.narg('method') OR sqlc.narg('method') IS NULL) AND
  (p.status = sqlc.narg('status') OR sqlc.narg('status') IS NULL) AND
  (p.date_created >= sqlc.narg('date_created_from') OR sqlc.narg('date_created_from') IS NULL) AND
  (p.date_created <= sqlc.narg('date_created_to') OR sqlc.narg('date_created_to') IS NULL)
)
ORDER BY p.date_created DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreatePayment :one
INSERT INTO "payment"."base" (account_id, method, status, total)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdatePayment :one
UPDATE "payment"."base"
SET
    method = COALESCE(sqlc.narg('method'), method),
    status = COALESCE(sqlc.narg('status'), status),
    total = COALESCE(sqlc.narg('total'), total)
WHERE id = $1
RETURNING *;

-- name: DeletePayment :exec
DELETE FROM "payment"."base"
WHERE id = $1;

-- name: CreatePaymentItem :one
INSERT INTO "payment"."item" (payment_id, name, price)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreatePaymentVnpay :one
INSERT INTO "payment"."vnpay" (id, "vnp_TxnRef", "vnp_OrderInfo", "vnp_TransactionNo", "vnp_TransactionDate", "vnp_CreateDate", "vnp_IpAddr")
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
