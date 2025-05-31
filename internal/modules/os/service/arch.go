package ossvc

import (
	"context"

	osmodel "github.com/wagecloud/wagecloud-server/internal/modules/os/model"
	osstorage "github.com/wagecloud/wagecloud-server/internal/modules/os/storage"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
)

func (s *ServiceImpl) GetArch(ctx context.Context, id string) (osmodel.Arch, error) {
	return s.storage.GetArch(ctx, id)
}

type ListArchsParams struct {
	pagination.PaginationParams
	ID            *string
	Name          *string
	CreatedAtFrom *int64
	CreatedAtTo   *int64
}

func (s *ServiceImpl) ListArchs(ctx context.Context, params ListArchsParams) (res pagination.PaginateResult[osmodel.Arch], err error) {
	total, err := s.storage.CountArchs(ctx, osstorage.ListArchsParams{
		PaginationParams: params.PaginationParams,
		ID:               params.ID,
		Name:             params.Name,
		CreatedAtFrom:    params.CreatedAtFrom,
		CreatedAtTo:      params.CreatedAtTo,
	})
	if err != nil {
		return res, err
	}

	archs, err := s.storage.ListArchs(ctx, osstorage.ListArchsParams{
		PaginationParams: params.PaginationParams,
		ID:               params.ID,
		Name:             params.Name,
		CreatedAtFrom:    params.CreatedAtFrom,
		CreatedAtTo:      params.CreatedAtTo,
	})
	if err != nil {
		return res, err
	}

	return pagination.PaginateResult[osmodel.Arch]{
		Total: total,
		Limit: params.Limit,
		Page:  params.Offset(),
		Data:  archs,
	}, nil
}

func (s *ServiceImpl) CreateArch(ctx context.Context, arch osmodel.Arch) (osmodel.Arch, error) {
	return s.storage.CreateArch(ctx, arch)
}

type UpdateArchParams struct {
	ID    string
	NewID *string
	Name  *string
}

func (s *ServiceImpl) UpdateArch(ctx context.Context, params UpdateArchParams) (osmodel.Arch, error) {
	return s.storage.UpdateArch(ctx, osstorage.UpdateArchParams{
		ID:    params.ID,
		NewID: params.NewID,
		Name:  params.Name,
	})
}

func (s *ServiceImpl) DeleteArch(ctx context.Context, id string) error {
	return s.storage.DeleteArch(ctx, id)
}
