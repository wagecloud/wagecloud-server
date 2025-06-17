package instancestorage

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wagecloud/wagecloud-server/gen/sqlc"
	instancemodel "github.com/wagecloud/wagecloud-server/internal/modules/instance/model"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	pgxptr "github.com/wagecloud/wagecloud-server/internal/utils/pgx/ptr"
)

func (r *Storage) GetDomain(ctx context.Context, id int64) (instancemodel.Domain, error) {
	domain, err := r.sqlc.GetDomain(ctx, id)
	if err != nil {
		return instancemodel.Domain{}, err
	}

	return instancemodel.Domain{
		ID:        domain.ID,
		NetworkID: domain.NetworkID,
		Name:      domain.Name,
	}, nil
}

type ListDomainsParams struct {
	pagination.PaginationParams
	NetworkID *int64
	Name      *string
}

func (r *Storage) CountDomains(ctx context.Context, params ListDomainsParams) (int64, error) {
	return r.sqlc.CountDomains(ctx, sqlc.CountDomainsParams{
		NetworkID: *pgxptr.PtrToPgtype(&pgtype.Int8{}, params.NetworkID),
		Name:      *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Name),
	})
}

func (r *Storage) ListDomains(ctx context.Context, params ListDomainsParams) ([]instancemodel.Domain, error) {
	domains, err := r.sqlc.ListDomains(ctx, sqlc.ListDomainsParams{
		Offset:    params.Offset(),
		Limit:     params.Limit,
		NetworkID: *pgxptr.PtrToPgtype(&pgtype.Int8{}, params.NetworkID),
		Name:      *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Name),
	})
	if err != nil {
		return nil, err
	}

	var result []instancemodel.Domain
	for _, domain := range domains {
		result = append(result, instancemodel.Domain{
			ID:        domain.ID,
			NetworkID: domain.NetworkID,
			Name:      domain.Name,
		})
	}

	return result, nil
}

func (r *Storage) CreateDomain(ctx context.Context, domain instancemodel.Domain) (instancemodel.Domain, error) {
	row, err := r.sqlc.CreateDomain(ctx, sqlc.CreateDomainParams{
		NetworkID: domain.NetworkID,
		Name:      domain.Name,
	})
	if err != nil {
		return instancemodel.Domain{}, err
	}

	return instancemodel.Domain{
		ID:        row.ID,
		NetworkID: row.NetworkID,
		Name:      row.Name,
	}, nil
}

type UpdateDomainParams struct {
	ID   int64
	Name *string
}

func (r *Storage) UpdateDomain(ctx context.Context, params UpdateDomainParams) (instancemodel.Domain, error) {
	row, err := r.sqlc.UpdateDomain(ctx, sqlc.UpdateDomainParams{
		ID:   params.ID,
		Name: *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Name),
	})
	if err != nil {
		return instancemodel.Domain{}, err
	}

	return instancemodel.Domain{
		ID:        row.ID,
		NetworkID: row.NetworkID,
		Name:      row.Name,
	}, nil
}

func (r *Storage) DeleteDomain(ctx context.Context, id int64) error {
	return r.sqlc.DeleteDomain(ctx, id)
}
