package network

import (
	"context"

	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/repository"
)

var _ ServiceInterface = (*Service)(nil)

type Service struct {
	repo repository.Repository
}

type ServiceInterface interface {
	GetNetwork(ctx context.Context, params GetNetworkParams) (model.Network, error)
	ListNetworks(ctx context.Context, params ListNetworksParams) (model.PaginateResult[model.Network], error)
	CreateNetwork(ctx context.Context, params CreateNetworkParams) (model.Network, error)
	UpdateNetwork(ctx context.Context, params UpdateNetworkParams) (model.Network, error)
	DeleteNetwork(ctx context.Context, params DeleteNetworkParams) error
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

type GetNetworkParams struct {
	ID string
}

func (s *Service) GetNetwork(ctx context.Context, params GetNetworkParams) (model.Network, error) {
	return s.repo.GetNetwork(ctx, params.ID)
}

type ListNetworksParams struct {
	model.PaginationParams
	ID            *string
	PrivateIP     *string
	CreatedAtFrom *int64
	CreatedAtTo   *int64
}

func (s *Service) ListNetworks(ctx context.Context, params ListNetworksParams) (res model.PaginateResult[model.Network], err error) {
	repoParams := repository.ListNetworksParams{
		PaginationParams: params.PaginationParams,
		ID:               params.ID,
		PrivateIP:        params.PrivateIP,
		CreatedAtFrom:    params.CreatedAtFrom,
		CreatedAtTo:      params.CreatedAtTo,
	}

	total, err := s.repo.CountNetworks(ctx, repoParams)
	if err != nil {
		return res, err
	}

	networks, err := s.repo.ListNetworks(ctx, repoParams)
	if err != nil {
		return res, err
	}

	return model.PaginateResult[model.Network]{
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

func (s *Service) CreateNetwork(ctx context.Context, params CreateNetworkParams) (model.Network, error) {
	return s.repo.CreateNetwork(ctx, model.Network{
		ID:        params.ID,
		PrivateIP: params.PrivateIP,
	})
}

type UpdateNetworkParams struct {
	ID        string
	NewID     *string
	PrivateIP *string
}

func (s *Service) UpdateNetwork(ctx context.Context, params UpdateNetworkParams) (model.Network, error) {
	return s.repo.UpdateNetwork(ctx, repository.UpdateNetworkParams{
		ID:        params.ID,
		NewID:     params.NewID,
		PrivateIP: params.PrivateIP,
	})
}

type DeleteNetworkParams struct {
	ID string
}

func (s *Service) DeleteNetwork(ctx context.Context, params DeleteNetworkParams) error {
	return s.repo.DeleteNetwork(ctx, params.ID)
}
