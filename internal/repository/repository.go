package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/wagecloud/wagecloud-server/gen/sqlc"
	pgxutil "github.com/wagecloud/wagecloud-server/internal/db/pgx"
	"github.com/wagecloud/wagecloud-server/internal/model"
)

var _ Repository = (*RepositoryImpl)(nil)

type Repository interface {
	Begin(ctx context.Context) (*RepositoryTx, error)

	// Account
	GetAccount(ctx context.Context, params GetAccountParams) (model.AccountBase, error)
	CountAccounts(ctx context.Context, params ListAccountsParams) (int64, error)
	ListAccounts(ctx context.Context, params ListAccountsParams) ([]model.AccountBase, error)
	CreateAccount(ctx context.Context, account model.AccountBase) (model.AccountBase, error)
	UpdateAccount(ctx context.Context, params UpdateAccountParams) (model.AccountBase, error)
	DeleteAccount(ctx context.Context, id int64) error

	// Arch
	GetArch(ctx context.Context, id string) (model.Arch, error)
	CountArchs(ctx context.Context, params ListArchsParams) (int64, error)
	ListArchs(ctx context.Context, params ListArchsParams) ([]model.Arch, error)
	CreateArch(ctx context.Context, arch model.Arch) (model.Arch, error)
	UpdateArch(ctx context.Context, params UpdateArchParams) (model.Arch, error)
	DeleteArch(ctx context.Context, id string) error

	// Network
	GetNetwork(ctx context.Context, id string) (model.Network, error)
	CountNetworks(ctx context.Context, params ListNetworksParams) (int64, error)
	ListNetworks(ctx context.Context, params ListNetworksParams) ([]model.Network, error)
	CreateNetwork(ctx context.Context, network model.Network) (model.Network, error)
	UpdateNetwork(ctx context.Context, params UpdateNetworkParams) (model.Network, error)
	DeleteNetwork(ctx context.Context, id string) error

	// OS
	GetOS(ctx context.Context, id string) (model.OS, error)
	CountOSs(ctx context.Context, params ListOSsParams) (int64, error)
	ListOSs(ctx context.Context, params ListOSsParams) ([]model.OS, error)
	CreateOS(ctx context.Context, os model.OS) (model.OS, error)
	UpdateOS(ctx context.Context, params UpdateOSParams) (model.OS, error)
	DeleteOS(ctx context.Context, id string) error

	// VM
	GetVM(ctx context.Context, params GetVMParams) (*model.VM, error)
	CountVMs(ctx context.Context, params ListVMParams) (int64, error)
	ListVMs(ctx context.Context, params ListVMParams) ([]*model.VM, error)
	CreateVM(ctx context.Context, vm model.VM) (model.VM, error)
	UpdateVM(ctx context.Context, params UpdateVMParams) (model.VM, error)
	DeleteVM(ctx context.Context, params DeleteVMParams) error
}

type RepositoryImpl struct {
	db   pgxutil.DBTX
	sqlc *sqlc.Queries
}

type RepositoryTx struct {
	*RepositoryImpl
	tx pgx.Tx
}

func NewRepository(db pgxutil.DBTX) *RepositoryImpl {
	return &RepositoryImpl{
		db:   db,
		sqlc: sqlc.New(db),
	}
}

func (r *RepositoryImpl) Begin(ctx context.Context) (*RepositoryTx, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &RepositoryTx{
		RepositoryImpl: NewRepository(tx),
		tx:             tx,
	}, nil
}

func (r *RepositoryTx) Commit(ctx context.Context) error {
	return r.tx.Commit(ctx)
}

func (r *RepositoryTx) Rollback(ctx context.Context) error {
	return r.tx.Rollback(ctx)
}
