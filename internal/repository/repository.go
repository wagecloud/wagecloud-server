package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/wagecloud/wagecloud-server/gen/sqlc"
	pgxutil "github.com/wagecloud/wagecloud-server/internal/db/pgx"
)

type Repository struct {
	db   pgxutil.DBTX
	sqlc *sqlc.Queries
}

type RepositoryTx struct {
	*Repository
	tx pgx.Tx
}

func NewRepository(db pgxutil.DBTX) *Repository {
	return &Repository{
		db:   db,
		sqlc: sqlc.New(db),
	}
}

func (r *Repository) Begin(ctx context.Context) (*RepositoryTx, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &RepositoryTx{
		Repository: NewRepository(tx),
		tx:         tx,
	}, nil
}

func (r *RepositoryTx) Commit(ctx context.Context) error {
	return r.tx.Commit(ctx)
}

func (r *RepositoryTx) Rollback(ctx context.Context) error {
	return r.tx.Rollback(ctx)
}
