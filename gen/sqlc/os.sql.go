// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: os.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countOSs = `-- name: CountOSs :one
SELECT COUNT(id)
FROM os
WHERE (
  (id ILIKE '%' || $1 || '%' OR $1 IS NULL) AND
  (name ILIKE '%' || $2 || '%' OR $2 IS NULL) AND
  (created_at >= $3 OR $3 IS NULL) AND
  (created_at <= $4 OR $4 IS NULL)
)
`

type CountOSsParams struct {
	ID            pgtype.Text
	Name          pgtype.Text
	CreatedAtFrom pgtype.Timestamptz
	CreatedAtTo   pgtype.Timestamptz
}

func (q *Queries) CountOSs(ctx context.Context, arg CountOSsParams) (int64, error) {
	row := q.db.QueryRow(ctx, countOSs,
		arg.ID,
		arg.Name,
		arg.CreatedAtFrom,
		arg.CreatedAtTo,
	)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createOS = `-- name: CreateOS :one
INSERT INTO os (id, name)
VALUES ($1, $2)
RETURNING id, name, created_at
`

type CreateOSParams struct {
	ID   string
	Name string
}

func (q *Queries) CreateOS(ctx context.Context, arg CreateOSParams) (O, error) {
	row := q.db.QueryRow(ctx, createOS, arg.ID, arg.Name)
	var i O
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return i, err
}

const deleteOS = `-- name: DeleteOS :exec
DELETE FROM os
WHERE id = $1
`

func (q *Queries) DeleteOS(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, deleteOS, id)
	return err
}

const getOS = `-- name: GetOS :one
SELECT os.id, os.name, os.created_at
FROM os
WHERE id = $1
`

func (q *Queries) GetOS(ctx context.Context, id string) (O, error) {
	row := q.db.QueryRow(ctx, getOS, id)
	var i O
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return i, err
}

const listOSs = `-- name: ListOSs :many
SELECT os.id, os.name, os.created_at
FROM os
WHERE (
  (id ILIKE '%' || $1 || '%' OR $1 IS NULL) AND
  (name ILIKE '%' || $2 || '%' OR $2 IS NULL) AND
  (created_at >= $3 OR $3 IS NULL) AND
  (created_at <= $4 OR $4 IS NULL)
)
ORDER BY created_at DESC
LIMIT $6
OFFSET $5
`

type ListOSsParams struct {
	ID            pgtype.Text
	Name          pgtype.Text
	CreatedAtFrom pgtype.Timestamptz
	CreatedAtTo   pgtype.Timestamptz
	Offset        int32
	Limit         int32
}

func (q *Queries) ListOSs(ctx context.Context, arg ListOSsParams) ([]O, error) {
	rows, err := q.db.Query(ctx, listOSs,
		arg.ID,
		arg.Name,
		arg.CreatedAtFrom,
		arg.CreatedAtTo,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []O
	for rows.Next() {
		var i O
		if err := rows.Scan(&i.ID, &i.Name, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateOS = `-- name: UpdateOS :one
UPDATE os
SET
    id = COALESCE($2, id),
    name = COALESCE($3, name)
WHERE id = $1
RETURNING id, name, created_at
`

type UpdateOSParams struct {
	ID    string
	NewID pgtype.Text
	Name  pgtype.Text
}

func (q *Queries) UpdateOS(ctx context.Context, arg UpdateOSParams) (O, error) {
	row := q.db.QueryRow(ctx, updateOS, arg.ID, arg.NewID, arg.Name)
	var i O
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return i, err
}
