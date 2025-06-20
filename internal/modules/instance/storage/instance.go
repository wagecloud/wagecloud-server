package instancestorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wagecloud/wagecloud-server/gen/sqlc"
	"github.com/wagecloud/wagecloud-server/internal/client/pgxpool"
	instancemodel "github.com/wagecloud/wagecloud-server/internal/modules/instance/model"
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

func (s *Storage) GetInstance(ctx context.Context, id string) (instancemodel.Instance, error) {
	row, err := s.sqlc.GetInstance(ctx, id)
	if err != nil {
		return instancemodel.Instance{}, err
	}

	return instancemodel.Instance{
		ID:        row.ID,
		AccountID: row.AccountID,
		OSID:      row.OsID,
		ArchID:    row.ArchID,
		RegionID:  row.RegionID,
		Name:      row.Name,
		CPU:       row.Cpu,
		RAM:       row.Ram,
		Storage:   row.Storage,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

type ListInstancesParams struct {
	pagination.PaginationParams
	AccountID     *int64
	Name          *string
	OsID          *string
	ArchID        *string
	RegionID      *string
	CpuFrom       *int64
	CpuTo         *int64
	RamFrom       *int64
	RamTo         *int64
	StorageFrom   *int64
	StorageTo     *int64
	CreatedAtFrom *int64
	CreatedAtTo   *int64
}

func (s *Storage) CountInstances(ctx context.Context, params ListInstancesParams) (int64, error) {
	return s.sqlc.CountInstances(ctx, sqlc.CountInstancesParams{
		AccountID:     *pgxptr.PtrToPgtype(&pgtype.Int8{}, params.AccountID),
		Name:          *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Name),
		OsID:          *pgxptr.PtrToPgtype(&pgtype.Text{}, params.OsID),
		ArchID:        *pgxptr.PtrToPgtype(&pgtype.Text{}, params.ArchID),
		RegionID:      *pgxptr.PtrToPgtype(&pgtype.Text{}, params.RegionID),
		CpuFrom:       *pgxptr.PtrToPgtype(&pgtype.Int4{}, params.CpuFrom),
		CpuTo:         *pgxptr.PtrToPgtype(&pgtype.Int4{}, params.CpuTo),
		RamFrom:       *pgxptr.PtrToPgtype(&pgtype.Int4{}, params.RamFrom),
		RamTo:         *pgxptr.PtrToPgtype(&pgtype.Int4{}, params.CpuTo),
		StorageFrom:   *pgxptr.PtrToPgtype(&pgtype.Int4{}, params.StorageFrom),
		StorageTo:     *pgxptr.PtrToPgtype(&pgtype.Int4{}, params.StorageTo),
		CreatedAtFrom: *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtFrom)),
		CreatedAtTo:   *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtTo)),
	})
}

func (s *Storage) ListInstances(ctx context.Context, params ListInstancesParams) ([]instancemodel.Instance, error) {
	rows, err := s.sqlc.ListInstances(ctx, sqlc.ListInstancesParams{
		Limit:         int32(params.Limit),
		Offset:        int32(params.Offset()),
		AccountID:     *pgxptr.PtrToPgtype(&pgtype.Int8{}, params.AccountID),
		Name:          *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Name),
		OsID:          *pgxptr.PtrToPgtype(&pgtype.Text{}, params.OsID),
		ArchID:        *pgxptr.PtrToPgtype(&pgtype.Text{}, params.ArchID),
		RegionID:      *pgxptr.PtrToPgtype(&pgtype.Text{}, params.RegionID),
		CpuFrom:       *pgxptr.PtrToPgtype(&pgtype.Int4{}, params.CpuFrom),
		CpuTo:         *pgxptr.PtrToPgtype(&pgtype.Int4{}, params.CpuTo),
		RamFrom:       *pgxptr.PtrToPgtype(&pgtype.Int4{}, params.RamFrom),
		RamTo:         *pgxptr.PtrToPgtype(&pgtype.Int4{}, params.RamTo),
		StorageFrom:   *pgxptr.PtrToPgtype(&pgtype.Int4{}, params.StorageFrom),
		StorageTo:     *pgxptr.PtrToPgtype(&pgtype.Int4{}, params.StorageTo),
		CreatedAtFrom: *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtFrom)),
		CreatedAtTo:   *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtTo)),
	})
	if err != nil {
		return nil, err
	}

	var instances []instancemodel.Instance
	for _, row := range rows {
		instances = append(instances, instancemodel.Instance{
			ID:        row.ID,
			AccountID: row.AccountID,
			OSID:      row.OsID,
			ArchID:    row.ArchID,
			RegionID:  row.RegionID,
			Name:      row.Name,
			CPU:       row.Cpu,
			RAM:       row.Ram,
			Storage:   row.Storage,
			CreatedAt: row.CreatedAt.Time,
		})
	}

	return instances, nil
}

func (s *Storage) CreateInstance(ctx context.Context, instance instancemodel.Instance) (instancemodel.Instance, error) {
	row, err := s.sqlc.CreateInstance(ctx, sqlc.CreateInstanceParams{
		ID:        instance.ID,
		AccountID: instance.AccountID,
		OsID:      instance.OSID,
		ArchID:    instance.ArchID,
		RegionID:  instance.RegionID,
		Name:      instance.Name,
		Cpu:       instance.CPU,
		Ram:       instance.RAM,
		Storage:   instance.Storage,
	})
	if err != nil {
		return instancemodel.Instance{}, err
	}

	return instancemodel.Instance{
		ID:        row.ID,
		AccountID: row.AccountID,
		OSID:      row.OsID,
		ArchID:    row.ArchID,
		RegionID:  row.RegionID,
		Name:      row.Name,
		CPU:       row.Cpu,
		RAM:       row.Ram,
		Storage:   row.Storage,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

type UpdateInstanceParams struct {
	ID      string
	Name    *string
	CPU     *int64
	RAM     *int64
	Storage *int64
}

func (s *Storage) UpdateInstance(ctx context.Context, params UpdateInstanceParams) (instancemodel.Instance, error) {
	row, err := s.sqlc.UpdateInstance(ctx, sqlc.UpdateInstanceParams{
		ID:      params.ID,
		Name:    *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Name),
		Cpu:     *pgxptr.PtrToPgtype(&pgtype.Int4{}, params.CPU),
		Ram:     *pgxptr.PtrToPgtype(&pgtype.Int4{}, params.RAM),
		Storage: *pgxptr.PtrToPgtype(&pgtype.Int4{}, params.Storage),
	})
	if err != nil {
		return instancemodel.Instance{}, err
	}

	return instancemodel.Instance{
		ID:        row.ID,
		AccountID: row.AccountID,
		OSID:      row.OsID,
		ArchID:    row.ArchID,
		RegionID:  row.RegionID,
		Name:      row.Name,
		CPU:       row.Cpu,
		RAM:       row.Ram,
		Storage:   row.Storage,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (s *Storage) DeleteInstance(ctx context.Context, id string) error {
	return s.sqlc.DeleteInstance(ctx, id)
}
