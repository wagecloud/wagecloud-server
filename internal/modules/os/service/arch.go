package ossvc

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
	GetArch(ctx context.Context, id string) (model.Arch, error)
	ListArchs(ctx context.Context, params ListArchsParams) (model.PaginateResult[model.Arch], error)
	CreateArch(ctx context.Context, arch model.Arch) (model.Arch, error)
	UpdateArch(ctx context.Context, params UpdateArchParams) (model.Arch, error)
	DeleteArch(ctx context.Context, id string) error
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetArch(ctx context.Context, id string) (model.Arch, error) {
	return s.repo.GetArch(ctx, id)
}

type ListArchsParams struct {
	model.PaginationParams
	ID            *string
	Name          *string
	CreatedAtFrom *int64
	CreatedAtTo   *int64
}

func (s *Service) ListArchs(ctx context.Context, params ListArchsParams) (res model.PaginateResult[model.Arch], err error) {
	total, err := s.repo.CountArchs(ctx, repository.ListArchsParams{
		PaginationParams: params.PaginationParams,
		ID:               params.ID,
		Name:             params.Name,
		CreatedAtFrom:    params.CreatedAtFrom,
		CreatedAtTo:      params.CreatedAtTo,
	})
	if err != nil {
		return res, err
	}

	archs, err := s.repo.ListArchs(ctx, repository.ListArchsParams{
		PaginationParams: params.PaginationParams,
		ID:               params.ID,
		Name:             params.Name,
		CreatedAtFrom:    params.CreatedAtFrom,
		CreatedAtTo:      params.CreatedAtTo,
	})
	if err != nil {
		return res, err
	}

	return model.PaginateResult[model.Arch]{
		Total: total,
		Limit: params.Limit,
		Page:  params.Offset(),
		Data:  archs,
	}, nil
}

func (s *Service) CreateArch(ctx context.Context, arch model.Arch) (model.Arch, error) {
	return s.repo.CreateArch(ctx, arch)
}

type UpdateArchParams struct {
	ID    string
	NewID *string
	Name  *string
}

func (s *Service) UpdateArch(ctx context.Context, params UpdateArchParams) (model.Arch, error) {
	return s.repo.UpdateArch(ctx, repository.UpdateArchParams{
		ID:    params.ID,
		NewID: params.NewID,
		Name:  params.Name,
	})
}

func (s *Service) DeleteArch(ctx context.Context, id string) error {
	return s.repo.DeleteArch(ctx, id)
}
