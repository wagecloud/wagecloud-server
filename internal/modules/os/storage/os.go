package osstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wagecloud/wagecloud-server/gen/sqlc"
	"github.com/wagecloud/wagecloud-server/internal/client/pgxpool"
	osmodel "github.com/wagecloud/wagecloud-server/internal/modules/os/model"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	pgxptr "github.com/wagecloud/wagecloud-server/internal/utils/pgx/ptr"
	"github.com/wagecloud/wagecloud-server/internal/utils/ptr"
)

type Storage struct {
	db   pgxpool.DBTX
	sqlc *sqlc.Queries
}

type TxStorage struct {
	*Storage
	tx pgx.Tx
}

func NewStorage(db pgxpool.DBTX) *Storage {
	return &Storage{
		db:   db,
		sqlc: sqlc.New(db),
	}
}

func (s *Storage) BeginTx(ctx context.Context) (*TxStorage, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &TxStorage{
		Storage: NewStorage(tx),
		tx:      tx,
	}, nil
}

func (ts *TxStorage) Commit(ctx context.Context) error {
	return ts.tx.Commit(ctx)
}

func (ts *TxStorage) Rollback(ctx context.Context) error {
	return ts.tx.Rollback(ctx)
}

func (r *Storage) GetOS(ctx context.Context, id string) (osmodel.OS, error) {
	os, err := r.sqlc.GetOS(ctx, id)
	if err != nil {
		return osmodel.OS{}, err
	}

	return osmodel.OS{
		ID:        os.ID,
		Name:      os.Name,
		CreatedAt: os.CreatedAt.Time.UnixMilli(),
	}, nil
}

type ListOSsParams struct {
	pagination.PaginationParams
	ID            *string
	Name          *string
	CreatedAtFrom *int64
	CreatedAtTo   *int64
}

func (r *Storage) CountOSs(ctx context.Context, params ListOSsParams) (int64, error) {
	return r.sqlc.CountOSs(ctx, sqlc.CountOSsParams{
		ID:            *pgxptr.PtrToPgtype(&pgtype.Text{}, params.ID),
		Name:          *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Name),
		CreatedAtFrom: *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtFrom)),
		CreatedAtTo:   *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtTo)),
	})
}

func (r *Storage) ListOSs(ctx context.Context, params ListOSsParams) ([]osmodel.OS, error) {
	oss, err := r.sqlc.ListOSs(ctx, sqlc.ListOSsParams{
		Offset:        params.Offset(),
		Limit:         params.Limit,
		ID:            *pgxptr.PtrToPgtype(&pgtype.Text{}, params.ID),
		Name:          *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Name),
		CreatedAtFrom: *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtFrom)),
		CreatedAtTo:   *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtTo)),
	})
	if err != nil {
		return nil, err
	}

	var result []osmodel.OS
	for _, os := range oss {
		result = append(result, osmodel.OS{
			ID:        os.ID,
			Name:      os.Name,
			CreatedAt: os.CreatedAt.Time.UnixMilli(),
		})
	}

	return result, nil
}

func (r *Storage) CreateOS(ctx context.Context, os osmodel.OS) (osmodel.OS, error) {
	osCreated, err := r.sqlc.CreateOS(ctx, sqlc.CreateOSParams{
		ID:   os.ID,
		Name: os.Name,
	})
	if err != nil {
		return osmodel.OS{}, err
	}

	return osmodel.OS{
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

func (r *Storage) UpdateOS(ctx context.Context, params UpdateOSParams) (osmodel.OS, error) {
	osUpdated, err := r.sqlc.UpdateOS(ctx, sqlc.UpdateOSParams{
		ID:    params.ID,
		NewID: *pgxptr.PtrToPgtype(&pgtype.Text{}, params.NewID),
		Name:  *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Name),
	})
	if err != nil {
		return osmodel.OS{}, err
	}

	return osmodel.OS{
		ID:        osUpdated.ID,
		Name:      osUpdated.Name,
		CreatedAt: osUpdated.CreatedAt.Time.UnixMilli(),
	}, nil
}

func (r *Storage) DeleteOS(ctx context.Context, id string) error {
	return r.sqlc.DeleteOS(ctx, id)
}
