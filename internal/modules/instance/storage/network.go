package instancestorage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wagecloud/wagecloud-server/gen/sqlc"
	instancemodel "github.com/wagecloud/wagecloud-server/internal/modules/instance/model"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	pgxptr "github.com/wagecloud/wagecloud-server/internal/utils/pgx/ptr"
)

type GetNetworkParams struct {
	ID         *int64
	InstanceID *string
}

func (r *Storage) GetNetwork(ctx context.Context, params GetNetworkParams) (instancemodel.Network, error) {
	if params.ID == nil && params.InstanceID == nil {
		return instancemodel.Network{}, errors.New("either ID or InstanceID must be provided")
	}

	network, err := r.sqlc.GetNetwork(ctx, sqlc.GetNetworkParams{
		ID:         *pgxptr.PtrToPgtype(&pgtype.Int8{}, params.ID),
		InstanceID: *pgxptr.PtrToPgtype(&pgtype.Text{}, params.InstanceID),
	})
	if err != nil {
		return instancemodel.Network{}, err
	}

	return instancemodel.Network{
		ID:         network.ID,
		InstanceID: network.InstanceID,
		PrivateIP:  network.PrivateIp,
		MacAddress: network.MacAddress,
		PublicIP:   pgxptr.PgtypeToPtr[string](network.PublicIp),
	}, nil
}

type ListNetworksParams struct {
	pagination.PaginationParams
	InstanceID *string
	PrivateIP  *string
	MacAddress *string
	PublicIP   *string
}

func (r *Storage) CountNetworks(ctx context.Context, params ListNetworksParams) (int64, error) {
	return r.sqlc.CountNetworks(ctx, sqlc.CountNetworksParams{
		InstanceID: *pgxptr.PtrToPgtype(&pgtype.Text{}, params.InstanceID),
		PrivateIp:  *pgxptr.PtrToPgtype(&pgtype.Text{}, params.PrivateIP),
		MacAddress: *pgxptr.PtrToPgtype(&pgtype.Text{}, params.MacAddress),
		PublicIp:   *pgxptr.PtrToPgtype(&pgtype.Text{}, params.PublicIP),
	})
}

func (r *Storage) ListNetworks(ctx context.Context, params ListNetworksParams) ([]instancemodel.Network, error) {
	networks, err := r.sqlc.ListNetworks(ctx, sqlc.ListNetworksParams{
		Offset:     params.Offset(),
		Limit:      params.Limit,
		InstanceID: *pgxptr.PtrToPgtype(&pgtype.Text{}, params.InstanceID),
		PrivateIp:  *pgxptr.PtrToPgtype(&pgtype.Text{}, params.PrivateIP),
		MacAddress: *pgxptr.PtrToPgtype(&pgtype.Text{}, params.MacAddress),
		PublicIp:   *pgxptr.PtrToPgtype(&pgtype.Text{}, params.PublicIP),
	})
	if err != nil {
		return nil, err
	}

	var result []instancemodel.Network
	for _, network := range networks {
		result = append(result, instancemodel.Network{
			ID:         network.ID,
			InstanceID: network.InstanceID,
			PrivateIP:  network.PrivateIp,
			MacAddress: network.MacAddress,
			PublicIP:   pgxptr.PgtypeToPtr[string](network.PublicIp),
		})
	}

	return result, nil
}

func (r *Storage) CreateNetwork(ctx context.Context, network instancemodel.Network) (instancemodel.Network, error) {
	row, err := r.sqlc.CreateNetwork(ctx, sqlc.CreateNetworkParams{
		InstanceID: network.InstanceID,
		PrivateIp:  network.PrivateIP,
		MacAddress: network.MacAddress,
		PublicIp:   *pgxptr.PtrToPgtype(&pgtype.Text{}, network.PublicIP),
	})
	if err != nil {
		return instancemodel.Network{}, err
	}

	return instancemodel.Network{
		ID:         row.ID,
		InstanceID: row.InstanceID,
		PrivateIP:  row.PrivateIp,
		MacAddress: row.MacAddress,
		PublicIP:   pgxptr.PgtypeToPtr[string](row.PublicIp),
	}, nil
}

type UpdateNetworkParams struct {
	ID           *int64
	InstanceID   *string
	PrivateIP    *string
	MacAddress   *string
	PublicIP     *string
	NullPublicIP bool
}

func (r *Storage) UpdateNetwork(ctx context.Context, params UpdateNetworkParams) (instancemodel.Network, error) {
	if params.ID == nil && params.InstanceID == nil {
		return instancemodel.Network{}, errors.New("either ID or InstanceID must be provided")
	}

	row, err := r.sqlc.UpdateNetwork(ctx, sqlc.UpdateNetworkParams{
		ID:           *pgxptr.PtrToPgtype(&pgtype.Int8{}, params.ID),
		InstanceID:   *pgxptr.PtrToPgtype(&pgtype.Text{}, params.InstanceID),
		PrivateIp:    *pgxptr.PtrToPgtype(&pgtype.Text{}, params.PrivateIP),
		MacAddress:   *pgxptr.PtrToPgtype(&pgtype.Text{}, params.MacAddress),
		PublicIp:     *pgxptr.PtrToPgtype(&pgtype.Text{}, params.PublicIP),
		NullPublicIp: params.NullPublicIP,
	})
	if err != nil {
		return instancemodel.Network{}, err
	}

	return instancemodel.Network{
		ID:         row.ID,
		InstanceID: row.InstanceID,
		PrivateIP:  row.PrivateIp,
		MacAddress: row.MacAddress,
		PublicIP:   pgxptr.PgtypeToPtr[string](row.PublicIp),
	}, nil
}

func (r *Storage) DeleteNetwork(ctx context.Context, id int64) error {
	return r.sqlc.DeleteNetwork(ctx, id)
}
