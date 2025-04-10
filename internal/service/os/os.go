package os

import (
	"context"
	"time"

	"github.com/google/uuid"
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
	id string
}

func (s *Service) GetOS(ctx context.Context, params GetOSParams) (model.OS, error) {
	os, err := s.repo.GetOS(ctx, params.id)

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
	CreatedAtFrom *int64
	CreatedAtTo   *int64
	Name          *string
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
	Name string
}

func (s *Service) CreateOS(ctx context.Context, params CreateOSParams) (model.OS, error) {
	osParams := model.OS{
		ID:        uuid.New().String(),
		Name:      params.Name,
		CreatedAt: time.Now().Local().UnixMicro(),
	}

	os, err := s.repo.CreateOS(ctx, osParams)

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
}

func (s *Service) UpdateOS(ctx context.Context, params UpdateOSParams) (model.OS, error) {
	return model.OS{}, nil
}

type DeleteOSParams struct {
}

func (s *Service) DeleteOS(ctx context.Context, params DeleteOSParams) error {
	return nil
}
