package instancesvc

import (
	"context"

	instancemodel "github.com/wagecloud/wagecloud-server/internal/modules/instance/model"
	instancestorage "github.com/wagecloud/wagecloud-server/internal/modules/instance/storage"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
)

func (s *ServiceImpl) GetRegion(ctx context.Context, id string) (instancemodel.Region, error) {
	return s.storage.GetRegion(ctx, id)
}

type ListRegionsParams struct {
	pagination.PaginationParams
	ID   *string
	Name *string
}

func (s *ServiceImpl) ListRegions(ctx context.Context, params ListRegionsParams) (res pagination.PaginateResult[instancemodel.Region], err error) {
	// TODO: rename all repoParams into storageParams
	storageParams := instancestorage.ListRegionsParams{
		PaginationParams: params.PaginationParams,
		ID:               params.ID,
		Name:             params.Name,
	}

	total, err := s.storage.CountRegions(ctx, storageParams)
	if err != nil {
		return res, err
	}

	regions, err := s.storage.ListRegions(ctx, storageParams)
	if err != nil {
		return res, err
	}

	return pagination.PaginateResult[instancemodel.Region]{
		Data:     regions,
		Limit:    params.Limit,
		Page:     params.Page,
		Total:    total,
		NextPage: params.NextPage(total),
	}, nil
}

type CreateRegionParams struct {
	ID   string
	Name string
}

func (s *ServiceImpl) CreateRegion(ctx context.Context, params CreateRegionParams) (instancemodel.Region, error) {
	return s.storage.CreateRegion(ctx, instancemodel.Region{
		ID:   params.ID,
		Name: params.Name,
	})
}

type UpdateRegionParams struct {
	ID    string
	NewID *string
	Name  *string
}

func (s *ServiceImpl) UpdateRegion(ctx context.Context, params UpdateRegionParams) (instancemodel.Region, error) {
	return s.storage.UpdateRegion(ctx, instancestorage.UpdateRegionParams{
		ID:    params.ID,
		NewID: params.NewID,
		Name:  params.Name,
	})
}

func (s *ServiceImpl) DeleteRegion(ctx context.Context, id string) error {
	return s.storage.DeleteRegion(ctx, id)
}
