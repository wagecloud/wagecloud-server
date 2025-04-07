package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wagecloud/wagecloud-server/gen/sqlc"
	"github.com/wagecloud/wagecloud-server/internal/model"
)

func (r *RepositoryImpl) GetAccountByID(ctx context.Context, accountID int64) (model.AccountBase, error) {
	baseRow, err := r.sqlc.GetAccountByID(ctx, accountID)
	if err != nil {
		return model.AccountBase{}, err
	}

	return model.AccountBase{
		ID:        baseRow.ID,
		Name:      baseRow.Name.String,
		Email:     baseRow.Email,
		CreatedAt: baseRow.CreatedAt.Time,
	}, nil
}

func (r *RepositoryImpl) CreateAccount(ctx context.Context, account model.AccountBase) (model.AccountBase, error) {
	createAccountParams := sqlc.CreateAccountParams{
		Name:  pgtype.Text{String: account.Name, Valid: true},
		Email: account.Email,
	}

	row, err := r.sqlc.CreateAccount(ctx, createAccountParams)

	if err != nil {
		return model.AccountBase{}, err
	}

	return model.AccountBase{
		ID:        row.ID,
		Name:      row.Name.String,
		Email:     row.Email,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (r *RepositoryImpl) DeleteAccount(ctx context.Context, accountID int64) error {
	return r.sqlc.DeleteAccount(ctx, accountID)
}
