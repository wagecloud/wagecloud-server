package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wagecloud/wagecloud-server/gen/sqlc"
	pgxutil "github.com/wagecloud/wagecloud-server/internal/db/pgx"
	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/util/ptr"
)

type GetVMParams struct {
	ID        string
	AccountID *int64
}

func (r *RepositoryImpl) GetVM(ctx context.Context, params GetVMParams) (model.VM, error) {
	row, err := r.sqlc.GetVM(ctx, sqlc.GetVMParams{
		ID:        params.ID,
		AccountID: *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AccountID),
	})
	if err != nil {
		return model.VM{}, err
	}

	return model.VM{
		ID:        row.ID,
		NetworkID: row.NetworkID,
		OsID:      row.OsID,
		ArchID:    row.ArchID,
		Name:      row.Name,
		Cpu:       row.Cpu,
		Ram:       row.Ram,
		Storage:   row.Storage,
	}, nil
}

type ListVMsParams struct {
	model.PaginationParams
	AccountID     *int64
	NetworkID     *string
	OsID          *string
	ArchID        *string
	Name          *string
	CpuFrom       *int64
	CpuTo         *int64
	RamFrom       *int64
	RamTo         *int64
	StorageFrom   *int64
	StorageTo     *int64
	CreatedAtFrom *int64
	CreatedAtTo   *int64
}

func (r *RepositoryImpl) CountVMs(ctx context.Context, params ListVMsParams) (int64, error) {
	return r.sqlc.CountVMs(ctx, sqlc.CountVMsParams{
		AccountID:     *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AccountID),
		NetworkID:     *pgxutil.PtrToPgtype(&pgtype.Text{}, params.NetworkID),
		OsID:          *pgxutil.PtrToPgtype(&pgtype.Text{}, params.OsID),
		ArchID:        *pgxutil.PtrToPgtype(&pgtype.Text{}, params.ArchID),
		Name:          *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name),
		CpuFrom:       *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.CpuFrom),
		CpuTo:         *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.CpuTo),
		RamFrom:       *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.RamFrom),
		RamTo:         *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.RamTo),
		StorageFrom:   *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.StorageFrom),
		StorageTo:     *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.StorageTo),
		CreatedAtFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtFrom)),
		CreatedAtTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtTo)),
	})
}

func (r *RepositoryImpl) ListVMs(ctx context.Context, params ListVMsParams) ([]model.VM, error) {
	rows, err := r.sqlc.ListVMs(ctx, sqlc.ListVMsParams{
		Offset:        params.Offset(),
		Limit:         params.Limit,
		AccountID:     *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AccountID),
		NetworkID:     *pgxutil.PtrToPgtype(&pgtype.Text{}, params.NetworkID),
		OsID:          *pgxutil.PtrToPgtype(&pgtype.Text{}, params.OsID),
		ArchID:        *pgxutil.PtrToPgtype(&pgtype.Text{}, params.ArchID),
		Name:          *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name),
		CpuFrom:       *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.CpuFrom),
		CpuTo:         *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.CpuTo),
		RamFrom:       *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.RamFrom),
		RamTo:         *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.RamTo),
		StorageFrom:   *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.StorageFrom),
		StorageTo:     *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.StorageTo),
		CreatedAtFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtFrom)),
		CreatedAtTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtTo)),
	})
	if err != nil {
		return nil, err
	}

	vms := make([]model.VM, 0, len(rows))
	for _, row := range rows {
		vms = append(vms, model.VM{
			ID:        row.ID,
			NetworkID: row.NetworkID,
			OsID:      row.OsID,
			ArchID:    row.ArchID,
			Name:      row.Name,
			Cpu:       row.Cpu,
			Ram:       row.Ram,
			Storage:   row.Storage,
		})
	}

	return vms, nil
}

func (r *RepositoryImpl) CreateVM(ctx context.Context, vm model.VM) (model.VM, error) {
	row, err := r.sqlc.CreateVM(ctx, sqlc.CreateVMParams{
		AccountID: vm.AccountID,
		NetworkID: vm.NetworkID,
		OsID:      vm.OsID,
		ArchID:    vm.ArchID,
		Name:      vm.Name,
		Cpu:       vm.Cpu,
		Ram:       vm.Ram,
		Storage:   vm.Storage,
	})
	if err != nil {
		return model.VM{}, err
	}

	return model.VM{
		ID:        row.ID,
		AccountID: row.AccountID,
		NetworkID: row.NetworkID,
		OsID:      row.OsID,
		ArchID:    row.ArchID,
		Name:      row.Name,
		Cpu:       row.Cpu,
		Ram:       row.Ram,
		Storage:   row.Storage,
	}, nil
}

type UpdateVMParams struct {
	ID        string
	AccountID *int64
	NetworkID *string
	OsID      *string
	ArchID    *string
	Name      *string
	Cpu       *int32
	Ram       *int32
	Storage   *int32
}

func (r *RepositoryImpl) UpdateVM(ctx context.Context, params UpdateVMParams) (model.VM, error) {
	row, err := r.sqlc.UpdateVM(ctx, sqlc.UpdateVMParams{
		ID:        params.ID,
		AccountID: *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AccountID),
		NetworkID: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.NetworkID),
		OsID:      *pgxutil.PtrToPgtype(&pgtype.Text{}, params.OsID),
		ArchID:    *pgxutil.PtrToPgtype(&pgtype.Text{}, params.ArchID),
		Name:      *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name),
		Cpu:       *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.Cpu),
		Ram:       *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.Ram),
		Storage:   *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.Storage),
	})
	if err != nil {
		return model.VM{}, err
	}

	return model.VM{
		ID:        row.ID,
		AccountID: row.AccountID,
		NetworkID: row.NetworkID,
		OsID:      row.OsID,
		ArchID:    row.ArchID,
		Name:      row.Name,
		Cpu:       row.Cpu,
		Ram:       row.Ram,
		Storage:   row.Storage,
	}, nil
}

type DeleteVMParams struct {
	ID        string
	AccountID *int64
}

func (r *RepositoryImpl) DeleteVM(ctx context.Context, params DeleteVMParams) error {
	return r.sqlc.DeleteVM(ctx, sqlc.DeleteVMParams{
		ID:        params.ID,
		AccountID: *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AccountID),
	})
}
