package ossvc

import (
	"context"

	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/repository"
	"github.com/wagecloud/wagecloud-server/internal/service/libvirt"
)

var _ ServiceInterface = (*Service)(nil)

type Service struct {
	repo    repository.Repository
	libvirt libvirt.ServiceInterface
}

type ServiceInterface interface {
	GetOS(ctx context.Context, params GetOSParams) (model.OS, error)
	ListOSs(ctx context.Context, params ListOSsParams) (model.PaginateResult[model.OS], error)
	CreateOS(ctx context.Context, params CreateOSParams) (model.OS, error)
	UpdateOS(ctx context.Context, params UpdateOSParams) (model.OS, error)
	DeleteOS(ctx context.Context, params DeleteOSParams) error
}

func NewService(repo repository.Repository, libvirt libvirt.ServiceInterface) *Service {
	return &Service{repo: repo, libvirt: libvirt}
}

type GetOSParams struct {
	ID string
}

func (s *Service) GetOS(ctx context.Context, params GetOSParams) (model.OS, error) {
	os, err := s.repo.GetOS(ctx, params.ID)

	if err != nil {
		return model.OS{}, err
	}

	return model.OS{
		ID:        os.ID,
		Name:      os.Name,
		CreatedAt: os.CreatedAt,
	}, nil
}

type ListOSsParams struct {
	model.PaginationParams
	ID            *string
	Name          *string
	CreatedAtFrom *int64
	CreatedAtTo   *int64
}

func (s *Service) ListOSs(ctx context.Context, params ListOSsParams) (res model.PaginateResult[model.OS], err error) {
	repoParams := repository.ListOSsParams{
		PaginationParams: params.PaginationParams,
		ID:               params.ID,
		Name:             params.Name,
		CreatedAtFrom:    params.CreatedAtFrom,
		CreatedAtTo:      params.CreatedAtTo,
	}

	total, err := s.repo.CountOSs(ctx, repoParams)
	if err != nil {
		return res, err
	}

	oss, err := s.repo.ListOSs(ctx, repoParams)
	if err != nil {
		return res, err
	}

	return model.PaginateResult[model.OS]{
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

func (s *Service) CreateOS(ctx context.Context, params CreateOSParams) (model.OS, error) {
	os, err := s.repo.CreateOS(ctx, model.OS{
		ID:   params.ID,
		Name: params.Name,
	})
	if err != nil {
		return model.OS{}, err
	}

	return model.OS{
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

func (s *Service) UpdateOS(ctx context.Context, params UpdateOSParams) (model.OS, error) {
	os, err := s.repo.UpdateOS(ctx, repository.UpdateOSParams{
		ID:    params.ID,
		NewID: params.NewID,
		Name:  params.Name,
	})
	if err != nil {
		return model.OS{}, err
	}

	return model.OS{
		ID:        os.ID,
		Name:      os.Name,
		CreatedAt: os.CreatedAt,
	}, nil
}

type DeleteOSParams struct {
	ID string
}

func (s *Service) DeleteOS(ctx context.Context, params DeleteOSParams) error {
	return s.repo.DeleteOS(ctx, params.ID)
}
