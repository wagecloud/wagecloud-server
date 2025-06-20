// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: network.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countNetworks = `-- name: CountNetworks :one
SELECT COUNT(id)
FROM "instance"."network"
WHERE (
  (instance_id = $1 OR $1 IS NULL) AND
  (private_ip ILIKE '%' || $2 || '%' OR $2 IS NULL) AND
  (mac_address ILIKE '%' || $3 || '%' OR $3 IS NULL) AND
  (public_ip ILIKE '%' || $4 || '%' OR $4 IS NULL)
)
`

type CountNetworksParams struct {
	InstanceID pgtype.Text
	PrivateIp  pgtype.Text
	MacAddress pgtype.Text
	PublicIp   pgtype.Text
}

func (q *Queries) CountNetworks(ctx context.Context, arg CountNetworksParams) (int64, error) {
	row := q.db.QueryRow(ctx, countNetworks,
		arg.InstanceID,
		arg.PrivateIp,
		arg.MacAddress,
		arg.PublicIp,
	)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createNetwork = `-- name: CreateNetwork :one
INSERT INTO "instance"."network" (instance_id, private_ip, mac_address, public_ip)
VALUES ($1, $2, $3, $4)
RETURNING id, instance_id, private_ip, mac_address, public_ip
`

type CreateNetworkParams struct {
	InstanceID string
	PrivateIp  string
	MacAddress string
	PublicIp   pgtype.Text
}

func (q *Queries) CreateNetwork(ctx context.Context, arg CreateNetworkParams) (InstanceNetwork, error) {
	row := q.db.QueryRow(ctx, createNetwork,
		arg.InstanceID,
		arg.PrivateIp,
		arg.MacAddress,
		arg.PublicIp,
	)
	var i InstanceNetwork
	err := row.Scan(
		&i.ID,
		&i.InstanceID,
		&i.PrivateIp,
		&i.MacAddress,
		&i.PublicIp,
	)
	return i, err
}

const deleteNetwork = `-- name: DeleteNetwork :exec
DELETE FROM "instance"."network"
WHERE id = $1
`

func (q *Queries) DeleteNetwork(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteNetwork, id)
	return err
}

const getNetwork = `-- name: GetNetwork :one
SELECT network.id, network.instance_id, network.private_ip, network.mac_address, network.public_ip
FROM "instance"."network" network
WHERE (
  id = $1 OR instance_id = $2
)
`

type GetNetworkParams struct {
	ID         pgtype.Int8
	InstanceID pgtype.Text
}

func (q *Queries) GetNetwork(ctx context.Context, arg GetNetworkParams) (InstanceNetwork, error) {
	row := q.db.QueryRow(ctx, getNetwork, arg.ID, arg.InstanceID)
	var i InstanceNetwork
	err := row.Scan(
		&i.ID,
		&i.InstanceID,
		&i.PrivateIp,
		&i.MacAddress,
		&i.PublicIp,
	)
	return i, err
}

const listNetworks = `-- name: ListNetworks :many
SELECT network.id, network.instance_id, network.private_ip, network.mac_address, network.public_ip
FROM "instance"."network" network
WHERE (
  (instance_id = $1 OR $1 IS NULL) AND
  (private_ip ILIKE '%' || $2 || '%' OR $2 IS NULL) AND
  (mac_address ILIKE '%' || $3 || '%' OR $3 IS NULL) AND
  (public_ip ILIKE '%' || $4 || '%' OR $4 IS NULL)
)
ORDER BY id DESC
LIMIT $6
OFFSET $5
`

type ListNetworksParams struct {
	InstanceID pgtype.Text
	PrivateIp  pgtype.Text
	MacAddress pgtype.Text
	PublicIp   pgtype.Text
	Offset     int32
	Limit      int32
}

func (q *Queries) ListNetworks(ctx context.Context, arg ListNetworksParams) ([]InstanceNetwork, error) {
	rows, err := q.db.Query(ctx, listNetworks,
		arg.InstanceID,
		arg.PrivateIp,
		arg.MacAddress,
		arg.PublicIp,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []InstanceNetwork
	for rows.Next() {
		var i InstanceNetwork
		if err := rows.Scan(
			&i.ID,
			&i.InstanceID,
			&i.PrivateIp,
			&i.MacAddress,
			&i.PublicIp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateNetwork = `-- name: UpdateNetwork :one
UPDATE "instance"."network"
SET
    private_ip = COALESCE($1, private_ip),
    mac_address = COALESCE($2, mac_address),
    public_ip = CASE
        WHEN $3::boolean THEN NULL
        ELSE COALESCE($4, public_ip)
    END
WHERE (
  (id = $5 OR $5 IS NULL) AND
  (instance_id = $6 OR $6 IS NULL)
)
RETURNING id, instance_id, private_ip, mac_address, public_ip
`

type UpdateNetworkParams struct {
	PrivateIp    pgtype.Text
	MacAddress   pgtype.Text
	NullPublicIp bool
	PublicIp     pgtype.Text
	ID           pgtype.Int8
	InstanceID   pgtype.Text
}

func (q *Queries) UpdateNetwork(ctx context.Context, arg UpdateNetworkParams) (InstanceNetwork, error) {
	row := q.db.QueryRow(ctx, updateNetwork,
		arg.PrivateIp,
		arg.MacAddress,
		arg.NullPublicIp,
		arg.PublicIp,
		arg.ID,
		arg.InstanceID,
	)
	var i InstanceNetwork
	err := row.Scan(
		&i.ID,
		&i.InstanceID,
		&i.PrivateIp,
		&i.MacAddress,
		&i.PublicIp,
	)
	return i, err
}
