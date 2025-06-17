package instancestorage

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wagecloud/wagecloud-server/gen/sqlc"
	instancemodel "github.com/wagecloud/wagecloud-server/internal/modules/instance/model"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	pgxptr "github.com/wagecloud/wagecloud-server/internal/utils/pgx/ptr"
)

func (r *Storage) GetRegion(ctx context.Context, id string) (instancemodel.Region, error) {
	region, err := r.sqlc.GetRegion(ctx, id)
	if err != nil {
		return instancemodel.Region{}, err
	}

	return instancemodel.Region{
		ID:   region.ID,
		Name: region.Name,
	}, nil
}

type ListRegionsParams struct {
	pagination.PaginationParams
	ID   *string
	Name *string
}

func (r *Storage) CountRegions(ctx context.Context, params ListRegionsParams) (int64, error) {
	return r.sqlc.CountRegions(ctx, sqlc.CountRegionsParams{
		ID:   *pgxptr.PtrToPgtype(&pgtype.Text{}, params.ID),
		Name: *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Name),
	})
}

func (r *Storage) ListRegions(ctx context.Context, params ListRegionsParams) ([]instancemodel.Region, error) {
	regions, err := r.sqlc.ListRegions(ctx, sqlc.ListRegionsParams{
		Offset: params.Offset(),
		Limit:  params.Limit,
		ID:     *pgxptr.PtrToPgtype(&pgtype.Text{}, params.ID),
		Name:   *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Name),
	})
	if err != nil {
		return nil, err
	}

	var result []instancemodel.Region
	for _, region := range regions {
		result = append(result, instancemodel.Region{
			ID:   region.ID,
			Name: region.Name,
		})
	}

	return result, nil
}

func (r *Storage) CreateRegion(ctx context.Context, region instancemodel.Region) (instancemodel.Region, error) {
	row, err := r.sqlc.CreateRegion(ctx, sqlc.CreateRegionParams{
		ID:   region.ID,
		Name: region.Name,
	})
	if err != nil {
		return instancemodel.Region{}, err
	}

	return instancemodel.Region{
		ID:   row.ID,
		Name: row.Name,
	}, nil
}

type UpdateRegionParams struct {
	ID    string
	NewID *string
	Name  *string
}

func (r *Storage) UpdateRegion(ctx context.Context, params UpdateRegionParams) (instancemodel.Region, error) {
	row, err := r.sqlc.UpdateRegion(ctx, sqlc.UpdateRegionParams{
		ID:    params.ID,
		NewID: *pgxptr.PtrToPgtype(&pgtype.Text{}, params.NewID),
		Name:  *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Name),
	})
	if err != nil {
		return instancemodel.Region{}, err
	}

	return instancemodel.Region{
		ID:   row.ID,
		Name: row.Name,
	}, nil
}

func (r *Storage) DeleteRegion(ctx context.Context, id string) error {
	return r.sqlc.DeleteRegion(ctx, id)
}
