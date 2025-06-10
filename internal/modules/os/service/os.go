package ossvc

import (
	"context"

	osmodel "github.com/wagecloud/wagecloud-server/internal/modules/os/model"
	osstorage "github.com/wagecloud/wagecloud-server/internal/modules/os/storage"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
)

type ServiceImpl struct {
	storage *osstorage.Storage
}

type Service interface {
	// OS
	GetOS(ctx context.Context, params GetOSParams) (osmodel.OS, error)
	ListOSs(ctx context.Context, params ListOSsParams) (pagination.PaginateResult[osmodel.OS], error)
	CreateOS(ctx context.Context, params CreateOSParams) (osmodel.OS, error)
	UpdateOS(ctx context.Context, params UpdateOSParams) (osmodel.OS, error)
	DeleteOS(ctx context.Context, params DeleteOSParams) error

	// Arch
	GetArch(ctx context.Context, id string) (osmodel.Arch, error)
	ListArchs(ctx context.Context, params ListArchsParams) (pagination.PaginateResult[osmodel.Arch], error)
	CreateArch(ctx context.Context, arch CreateArchParams) (osmodel.Arch, error)
	UpdateArch(ctx context.Context, params UpdateArchParams) (osmodel.Arch, error)
	DeleteArch(ctx context.Context, id string) error
}

func NewService(storage *osstorage.Storage) Service {
	return &ServiceImpl{
		storage: storage,
	}
}

type GetOSParams struct {
	ID string
}

func (s *ServiceImpl) GetOS(ctx context.Context, params GetOSParams) (osmodel.OS, error) {
	os, err := s.storage.GetOS(ctx, params.ID)

	if err != nil {
		return osmodel.OS{}, err
	}

	return osmodel.OS{
		ID:        os.ID,
		Name:      os.Name,
		CreatedAt: os.CreatedAt,
	}, nil
}

type ListOSsParams struct {
	pagination.PaginationParams
	ID            *string
	Name          *string
	CreatedAtFrom *int64
	CreatedAtTo   *int64
}

func (s *ServiceImpl) ListOSs(ctx context.Context, params ListOSsParams) (res pagination.PaginateResult[osmodel.OS], err error) {
	repoParams := osstorage.ListOSsParams{
		PaginationParams: params.PaginationParams,
		ID:               params.ID,
		Name:             params.Name,
		CreatedAtFrom:    params.CreatedAtFrom,
		CreatedAtTo:      params.CreatedAtTo,
	}

	total, err := s.storage.CountOSs(ctx, repoParams)
	if err != nil {
		return res, err
	}

	oss, err := s.storage.ListOSs(ctx, repoParams)
	if err != nil {
		return res, err
	}

	return pagination.PaginateResult[osmodel.OS]{
		Total: total,
		Page:  params.Page,
		Limit: params.Limit,
		Data:  oss,
	}, nil
}

type CreateOSParams struct {
	ID   string
	Name string
}

func (s *ServiceImpl) CreateOS(ctx context.Context, params CreateOSParams) (osmodel.OS, error) {
	os, err := s.storage.CreateOS(ctx, osmodel.OS{
		ID:   params.ID,
		Name: params.Name,
	})
	if err != nil {
		return osmodel.OS{}, err
	}

	return osmodel.OS{
		ID:        os.ID,
		Name:      os.Name,
		CreatedAt: os.CreatedAt,
	}, nil
}

type UpdateOSParams struct {
	ID    string
	NewID *string
	Name  *string
}

func (s *ServiceImpl) UpdateOS(ctx context.Context, params UpdateOSParams) (osmodel.OS, error) {
	os, err := s.storage.UpdateOS(ctx, osstorage.UpdateOSParams{
		ID:    params.ID,
		NewID: params.NewID,
		Name:  params.Name,
	})
	if err != nil {
		return osmodel.OS{}, err
	}

	return osmodel.OS{
		ID:        os.ID,
		Name:      os.Name,
		CreatedAt: os.CreatedAt,
	}, nil
}

type DeleteOSParams struct {
	ID string
}

func (s *ServiceImpl) DeleteOS(ctx context.Context, params DeleteOSParams) error {
	return s.storage.DeleteOS(ctx, params.ID)
}
