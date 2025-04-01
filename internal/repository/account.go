package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wagecloud/wagecloud-server/gen/sqlc"
	"github.com/wagecloud/wagecloud-server/internal/model"
)

func (r *Repository) GetAccountByID(ctx context.Context, accountID int64) (*model.AccountBase, error) {
	baseRow, err := r.sqlc.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	return &model.AccountBase{
		ID:        baseRow.ID,
		Name:      baseRow.Name.String,
		Email:     baseRow.Email,
		CreatedAt: baseRow.CreatedAt.Time,
	}, nil
}

func (r *Repository) CreateAccount(ctx context.Context, account *model.AccountBase) (*model.AccountBase, error) {
	createAccountParams := sqlc.CreateAccountParams{
		Name:  pgtype.Text{String: account.Name, Valid: true},
		Email: account.Email,
	}

	row, err := r.sqlc.CreateAccount(ctx, createAccountParams)

	if err != nil {
		return nil, err
	}

	return &model.AccountBase{
		ID:        row.ID,
		Name:      row.Name.String,
		Email:     row.Email,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (r *Repository) DeleteAccount(ctx context.Context, accountID int64) error {
	return r.sqlc.DeleteAccount(ctx, accountID)
}
