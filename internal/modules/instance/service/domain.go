package instancesvc

import (
	"context"

	instancemodel "github.com/wagecloud/wagecloud-server/internal/modules/instance/model"
	instancestorage "github.com/wagecloud/wagecloud-server/internal/modules/instance/storage"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
)

func (s *ServiceImpl) GetDomain(ctx context.Context, id int64) (instancemodel.Domain, error) {
	return s.storage.GetDomain(ctx, id)
}

type ListDomainsParams struct {
	pagination.PaginationParams
	NetworkID *int64
	Name      *string
}

func (s *ServiceImpl) ListDomains(ctx context.Context, params ListDomainsParams) (res pagination.PaginateResult[instancemodel.Domain], err error) {
	// TODO: rename all repoParams into storageParams
	storageParams := instancestorage.ListDomainsParams{
		PaginationParams: params.PaginationParams,
		NetworkID:        params.NetworkID,
		Name:             params.Name,
	}

	total, err := s.storage.CountDomains(ctx, storageParams)
	if err != nil {
		return res, err
	}

	domains, err := s.storage.ListDomains(ctx, storageParams)
	if err != nil {
		return res, err
	}

	return pagination.PaginateResult[instancemodel.Domain]{
		Data:     domains,
		Limit:    params.Limit,
		Page:     params.Page,
		Total:    total,
		NextPage: params.NextPage(total),
	}, nil
}

type CreateDomainParams struct {
	ID        int64
	NetworkID int64
	Name      string
}

func (s *ServiceImpl) CreateDomain(ctx context.Context, params CreateDomainParams) (instancemodel.Domain, error) {
	return s.storage.CreateDomain(ctx, instancemodel.Domain{
		ID:        params.ID,
		NetworkID: params.NetworkID,
		Name:      params.Name,
	})
}

type UpdateDomainParams struct {
	ID   int64
	Name *string
}

func (s *ServiceImpl) UpdateDomain(ctx context.Context, params UpdateDomainParams) (instancemodel.Domain, error) {
	return s.storage.UpdateDomain(ctx, instancestorage.UpdateDomainParams{
		ID:   params.ID,
		Name: params.Name,
	})
}

func (s *ServiceImpl) DeleteDomain(ctx context.Context, id int64) error {
	return s.storage.DeleteDomain(ctx, id)
}
