package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wagecloud/wagecloud-server/gen/sqlc"
	pgxutil "github.com/wagecloud/wagecloud-server/internal/db/pgx"
	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/util/ptr"
)

func (r *RepositoryImpl) GetNetwork(ctx context.Context, id string) (model.Network, error) {
	network, err := r.sqlc.GetNetwork(ctx, id)
	if err != nil {
		return model.Network{}, err
	}

	return model.Network{
		ID:        network.ID,
		PrivateIP: network.PrivateIp,
		CreatedAt: network.CreatedAt.Time.UnixMilli(),
	}, nil
}

type ListNetworksParams struct {
	model.PaginationParams
	ID            *string
	PrivateIP     *string
	CreatedAtFrom *int64
	CreatedAtTo   *int64
}

func (r *RepositoryImpl) CountNetworks(ctx context.Context, params ListNetworksParams) (int64, error) {
	return r.sqlc.CountNetworks(ctx, sqlc.CountNetworksParams{
		ID:            *pgxutil.PtrToPgtype(&pgtype.Text{}, params.ID),
		PrivateIp:     *pgxutil.PtrToPgtype(&pgtype.Text{}, params.PrivateIP),
		CreatedAtFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtFrom)),
		CreatedAtTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtTo)),
	})
}

func (r *RepositoryImpl) ListNetworks(ctx context.Context, params ListNetworksParams) ([]model.Network, error) {
	networks, err := r.sqlc.ListNetworks(ctx, sqlc.ListNetworksParams{
		Offset:        params.Offset(),
		Limit:         params.Limit,
		ID:            *pgxutil.PtrToPgtype(&pgtype.Text{}, params.ID),
		PrivateIp:     *pgxutil.PtrToPgtype(&pgtype.Text{}, params.PrivateIP),
		CreatedAtFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtFrom)),
		CreatedAtTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtTo)),
	})
	if err != nil {
		return nil, err
	}

	var result []model.Network
	for _, network := range networks {
		result = append(result, model.Network{
			ID:        network.ID,
			PrivateIP: network.PrivateIp,
			CreatedAt: network.CreatedAt.Time.UnixMilli(),
		})
	}

	return result, nil
}

func (r *RepositoryImpl) CreateNetwork(ctx context.Context, network model.Network) (model.Network, error) {
	row, err := r.sqlc.CreateNetwork(ctx, sqlc.CreateNetworkParams{
		ID:        network.ID,
		PrivateIp: network.PrivateIP,
	})
	if err != nil {
		return model.Network{}, err
	}

	return model.Network{
		ID:        row.ID,
		PrivateIP: row.PrivateIp,
	}, nil
}

type UpdateNetworkParams struct {
	ID        string
	NewID     *string
	PrivateIP *string
}

func (r *RepositoryImpl) UpdateNetwork(ctx context.Context, params UpdateNetworkParams) (model.Network, error) {
	row, err := r.sqlc.UpdateNetwork(ctx, sqlc.UpdateNetworkParams{
		ID:        params.ID,
		NewID:     *pgxutil.PtrToPgtype(&pgtype.Text{}, params.NewID),
		PrivateIp: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.PrivateIP),
	})
	if err != nil {
		return model.Network{}, err
	}

	return model.Network{
		ID:        row.ID,
		PrivateIP: row.PrivateIp,
	}, nil
}

func (r *RepositoryImpl) DeleteNetwork(ctx context.Context, id string) error {
	return r.sqlc.DeleteNetwork(ctx, id)
}
