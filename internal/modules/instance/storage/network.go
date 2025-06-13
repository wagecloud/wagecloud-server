package instancestorage

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wagecloud/wagecloud-server/gen/sqlc"
	instancemodel "github.com/wagecloud/wagecloud-server/internal/modules/instance/model"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	pgxptr "github.com/wagecloud/wagecloud-server/internal/utils/pgx/ptr"
	"github.com/wagecloud/wagecloud-server/internal/utils/ptr"
)

func (r *Storage) GetNetwork(ctx context.Context, id string) (instancemodel.Network, error) {
	network, err := r.sqlc.GetNetwork(ctx, id)
	if err != nil {
		return instancemodel.Network{}, err
	}

	return instancemodel.Network{
		ID:        network.ID,
		PrivateIP: network.PrivateIp,
		CreatedAt: network.CreatedAt.Time.UnixMilli(),
	}, nil
}

type ListNetworksParams struct {
	pagination.PaginationParams
	ID            *string
	PrivateIP     *string
	CreatedAtFrom *int64
	CreatedAtTo   *int64
}

func (r *Storage) CountNetworks(ctx context.Context, params ListNetworksParams) (int64, error) {
	return r.sqlc.CountNetworks(ctx, sqlc.CountNetworksParams{
		ID:            *pgxptr.PtrToPgtype(&pgtype.Text{}, params.ID),
		PrivateIp:     *pgxptr.PtrToPgtype(&pgtype.Text{}, params.PrivateIP),
		CreatedAtFrom: *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtFrom)),
		CreatedAtTo:   *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtTo)),
	})
}

func (r *Storage) ListNetworks(ctx context.Context, params ListNetworksParams) ([]instancemodel.Network, error) {
	networks, err := r.sqlc.ListNetworks(ctx, sqlc.ListNetworksParams{
		Offset:        params.Offset(),
		Limit:         params.Limit,
		ID:            *pgxptr.PtrToPgtype(&pgtype.Text{}, params.ID),
		PrivateIp:     *pgxptr.PtrToPgtype(&pgtype.Text{}, params.PrivateIP),
		CreatedAtFrom: *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtFrom)),
		CreatedAtTo:   *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtTo)),
	})
	if err != nil {
		return nil, err
	}

	var result []instancemodel.Network
	for _, network := range networks {
		result = append(result, instancemodel.Network{
			ID:        network.ID,
			PrivateIP: network.PrivateIp,
			CreatedAt: network.CreatedAt.Time.UnixMilli(),
		})
	}

	return result, nil
}

func (r *Storage) CreateNetwork(ctx context.Context, network instancemodel.Network) (instancemodel.Network, error) {
	row, err := r.sqlc.CreateNetwork(ctx, sqlc.CreateNetworkParams{
		ID:        network.ID,
		PrivateIp: network.PrivateIP,
	})
	if err != nil {
		return instancemodel.Network{}, err
	}

	return instancemodel.Network{
		ID:        row.ID,
		PrivateIP: row.PrivateIp,
	}, nil
}

type UpdateNetworkParams struct {
	ID        string
	PrivateIP *string
}

func (r *Storage) UpdateNetwork(ctx context.Context, params UpdateNetworkParams) (instancemodel.Network, error) {
	row, err := r.sqlc.UpdateNetwork(ctx, sqlc.UpdateNetworkParams{
		ID:        params.ID,
		PrivateIp: *pgxptr.PtrToPgtype(&pgtype.Text{}, params.PrivateIP),
	})
	if err != nil {
		return instancemodel.Network{}, err
	}

	return instancemodel.Network{
		ID:        row.ID,
		PrivateIP: row.PrivateIp,
	}, nil
}

func (r *Storage) DeleteNetwork(ctx context.Context, id string) error {
	return r.sqlc.DeleteNetwork(ctx, id)
}
