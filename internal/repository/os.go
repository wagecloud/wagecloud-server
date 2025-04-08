package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wagecloud/wagecloud-server/gen/sqlc"
	pgxutil "github.com/wagecloud/wagecloud-server/internal/db/pgx"
	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/util/ptr"
)

func (r *RepositoryImpl) GetOS(ctx context.Context, id string) (model.OS, error) {
	os, err := r.sqlc.GetOS(ctx, id)
	if err != nil {
		return model.OS{}, err
	}

	return model.OS{
		ID:        os.ID,
		Name:      os.Name,
		CreatedAt: os.CreatedAt.Time.UnixMilli(),
	}, nil
}

type ListOSsParams struct {
	model.PaginationParams
	ID            *string
	Name          *string
	CreatedAtFrom *int64
	CreatedAtTo   *int64
}

func (r *RepositoryImpl) CountOSs(ctx context.Context, params ListOSsParams) (int64, error) {
	return r.sqlc.CountOSs(ctx, sqlc.CountOSsParams{
		ID:            *pgxutil.PtrToPgtype(&pgtype.Text{}, params.ID),
		Name:          *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name),
		CreatedAtFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtFrom)),
		CreatedAtTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtTo)),
	})
}

func (r *RepositoryImpl) ListOSs(ctx context.Context, params ListOSsParams) ([]model.OS, error) {
	oss, err := r.sqlc.ListOSs(ctx, sqlc.ListOSsParams{
		Offset:        params.Offset(),
		Limit:         params.Limit,
		ID:            *pgxutil.PtrToPgtype(&pgtype.Text{}, params.ID),
		Name:          *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name),
		CreatedAtFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtFrom)),
		CreatedAtTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtTo)),
	})
	if err != nil {
		return nil, err
	}

	var result []model.OS
	for _, os := range oss {
		result = append(result, model.OS{
			ID:        os.ID,
			Name:      os.Name,
			CreatedAt: os.CreatedAt.Time.UnixMilli(),
		})
	}

	return result, nil
}

func (r *RepositoryImpl) CreateOS(ctx context.Context, os model.OS) (model.OS, error) {
	osCreated, err := r.sqlc.CreateOS(ctx, os.Name)
	if err != nil {
		return model.OS{}, err
	}

	return model.OS{
		ID:        osCreated.ID,
		Name:      osCreated.Name,
		CreatedAt: osCreated.CreatedAt.Time.UnixMilli(),
	}, nil
}

type UpdateOSParams struct {
	ID    string
	NewID *string
	Name  *string
}

func (r *RepositoryImpl) UpdateOS(ctx context.Context, params UpdateOSParams) (model.OS, error) {
	osUpdated, err := r.sqlc.UpdateOS(ctx, sqlc.UpdateOSParams{
		ID:    params.ID,
		NewID: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.NewID),
		Name:  *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name),
	})
	if err != nil {
		return model.OS{}, err
	}

	return model.OS{
		ID:        osUpdated.ID,
		Name:      osUpdated.Name,
		CreatedAt: osUpdated.CreatedAt.Time.UnixMilli(),
	}, nil
}

func (r *RepositoryImpl) DeleteOS(ctx context.Context, id string) error {
	return r.sqlc.DeleteOS(ctx, id)
}
