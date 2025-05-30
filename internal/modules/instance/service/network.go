package instancesvc

import (
	"context"

	instancemodel "github.com/wagecloud/wagecloud-server/internal/modules/instance/model"
	instancestorage "github.com/wagecloud/wagecloud-server/internal/modules/instance/storage"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
)

type GetNetworkParams struct {
	ID string
}

func (s *ServiceImpl) GetNetwork(ctx context.Context, params GetNetworkParams) (instancemodel.Network, error) {
	return s.storage.GetNetwork(ctx, params.ID)
}

type ListNetworksParams struct {
	pagination.PaginationParams
	ID            *string
	PrivateIP     *string
	CreatedAtFrom *int64
	CreatedAtTo   *int64
}

func (s *ServiceImpl) ListNetworks(ctx context.Context, params ListNetworksParams) (res pagination.PaginateResult[instancemodel.Network], err error) {
	repoParams := instancestorage.ListNetworksParams{
		PaginationParams: params.PaginationParams,
		ID:               params.ID,
		PrivateIP:        params.PrivateIP,
		CreatedAtFrom:    params.CreatedAtFrom,
		CreatedAtTo:      params.CreatedAtTo,
	}

	total, err := s.storage.CountNetworks(ctx, repoParams)
	if err != nil {
		return res, err
	}

	networks, err := s.storage.ListNetworks(ctx, repoParams)
	if err != nil {
		return res, err
	}

	return pagination.PaginateResult[instancemodel.Network]{
		Data:     networks,
		Limit:    params.Limit,
		Page:     params.Page,
		Total:    total,
		NextPage: params.NextPage(total),
	}, nil
}

type CreateNetworkParams struct {
	ID        string
	PrivateIP string
}

func (s *ServiceImpl) CreateNetwork(ctx context.Context, params CreateNetworkParams) (instancemodel.Network, error) {
	return s.storage.CreateNetwork(ctx, instancemodel.Network{
		ID:        params.ID,
		PrivateIP: params.PrivateIP,
	})
}

type UpdateNetworkParams struct {
	ID        string
	NewID     *string
	PrivateIP *string
}

func (s *ServiceImpl) UpdateNetwork(ctx context.Context, params UpdateNetworkParams) (instancemodel.Network, error) {
	return s.storage.UpdateNetwork(ctx, instancestorage.UpdateNetworkParams{
		ID:        params.ID,
		NewID:     params.NewID,
		PrivateIP: params.PrivateIP,
	})
}

type DeleteNetworkParams struct {
	ID string
}

func (s *ServiceImpl) DeleteNetwork(ctx context.Context, params DeleteNetworkParams) error {
	return s.storage.DeleteNetwork(ctx, params.ID)
}
