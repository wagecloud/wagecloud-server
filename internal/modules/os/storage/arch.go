package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wagecloud/wagecloud-server/gen/sqlc"
	"github.com/wagecloud/wagecloud-server/internal/model"
	pgxptr "github.com/wagecloud/wagecloud-server/internal/utils/pgx/ptr"
	"github.com/wagecloud/wagecloud-server/internal/utils/ptr"
)

func (r *RepositoryImpl) GetArch(ctx context.Context, id string) (model.Arch, error) {
	row, err := r.sqlc.GetArch(ctx, id)
	if err != nil {
		return model.Arch{}, err
	}

	return model.Arch{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt.Time.UnixMilli(),
	}, nil
}

type ListArchsParams struct {
	model.PaginationParams
	ID            *string
	Name          *string
	CreatedAtFrom *int64
	CreatedAtTo   *int64
}

func (r *RepositoryImpl) CountArchs(ctx context.Context, params ListArchsParams) (int64, error) {
	return r.sqlc.CountArchs(ctx, sqlc.CountArchsParams{
		ID:            *pgxptr.PtrToPgtype(&pgtype.Text{}, params.ID),
		Name:          *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Name),
		CreatedAtFrom: *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtFrom)),
		CreatedAtTo:   *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtTo)),
	})
}

func (r *RepositoryImpl) ListArchs(ctx context.Context, params ListArchsParams) ([]model.Arch, error) {
	rows, err := r.sqlc.ListArchs(ctx, sqlc.ListArchsParams{
		Limit:         params.Limit,
		Offset:        params.Offset(),
		ID:            *pgxptr.PtrToPgtype(&pgtype.Text{}, params.ID),
		Name:          *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Name),
		CreatedAtFrom: *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtFrom)),
		CreatedAtTo:   *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtTo)),
	})
	if err != nil {
		return nil, err
	}

	var archs []model.Arch
	for _, row := range rows {
		archs = append(archs, model.Arch{
			ID:        row.ID,
			Name:      row.Name,
			CreatedAt: row.CreatedAt.Time.UnixMilli(),
		})
	}

	return archs, nil
}

func (r *RepositoryImpl) CreateArch(ctx context.Context, arch model.Arch) (model.Arch, error) {
	row, err := r.sqlc.CreateArch(ctx, sqlc.CreateArchParams{
		ID:   arch.ID,
		Name: arch.Name,
	})
	if err != nil {
		return model.Arch{}, err
	}

	return model.Arch{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt.Time.UnixMilli(),
	}, nil
}

type UpdateArchParams struct {
	ID    string
	NewID *string
	Name  *string
}

func (r *RepositoryImpl) UpdateArch(ctx context.Context, params UpdateArchParams) (model.Arch, error) {
	row, err := r.sqlc.UpdateArch(ctx, sqlc.UpdateArchParams{
		ID:    params.ID,
		NewID: *pgxptr.PtrToPgtype(&pgtype.Text{}, params.NewID),
		Name:  *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Name),
	})
	if err != nil {
		return model.Arch{}, err
	}

	return model.Arch{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt.Time.UnixMilli(),
	}, nil
}

func (r *RepositoryImpl) DeleteArch(ctx context.Context, id string) error {
	return r.sqlc.DeleteArch(ctx, id)
}
