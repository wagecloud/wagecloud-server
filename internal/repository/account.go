package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wagecloud/wagecloud-server/gen/sqlc"
	pgxutil "github.com/wagecloud/wagecloud-server/internal/db/pgx"
	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/util/ptr"
)

type GetAccountParams struct {
	ID       *int64
	Username *string
	Email    *string
}

func (r *RepositoryImpl) GetAccount(ctx context.Context, params GetAccountParams) (model.AccountBase, error) {
	if params.ID == nil && params.Username == nil && params.Email == nil {
		return model.AccountBase{}, fmt.Errorf("at least one of ID, Username, or Email must be provided")
	}

	row, err := r.sqlc.GetAccount(ctx, sqlc.GetAccountParams{
		ID:       *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.ID),
		Username: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Username),
		Email:    *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Email),
	})
	if err != nil {
		return model.AccountBase{}, err
	}

	return model.AccountBase{
		ID:        row.ID,
		Username:  row.Username,
		Email:     row.Email,
		Name:      row.Name,
		Password:  row.Password,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

type ListAccountsParams struct {
	model.PaginationParams
	ID            *string
	Username      *string
	Email         *string
	Name          *string
	CreatedAtFrom *int64
	CreatedAtTo   *int64
}

func (r *RepositoryImpl) CountAccounts(ctx context.Context, params ListAccountsParams) (int64, error) {
	return r.sqlc.CountAccounts(ctx, sqlc.CountAccountsParams{
		ID:            *pgxutil.PtrToPgtype(&pgtype.Text{}, params.ID),
		Username:      *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Username),
		Email:         *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Email),
		Name:          *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name),
		CreatedAtFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtFrom)),
		CreatedAtTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtTo)),
	})
}

func (r *RepositoryImpl) ListAccounts(ctx context.Context, params ListAccountsParams) ([]model.AccountBase, error) {
	rows, err := r.sqlc.ListAccounts(ctx, sqlc.ListAccountsParams{
		Limit:         params.Limit,
		Offset:        params.Offset(),
		ID:            *pgxutil.PtrToPgtype(&pgtype.Text{}, params.ID),
		Username:      *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Username),
		Email:         *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Email),
		Name:          *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name),
		CreatedAtFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtFrom)),
		CreatedAtTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtTo)),
	})
	if err != nil {
		return nil, err
	}

	var accounts []model.AccountBase
	for _, row := range rows {
		accounts = append(accounts, model.AccountBase{
			ID:        row.ID,
			Username:  row.Username,
			Name:      row.Name,
			Email:     row.Email,
			Password:  row.Password,
			CreatedAt: row.CreatedAt.Time,
		})
	}

	return accounts, nil
}

func (r *RepositoryImpl) CreateAccount(ctx context.Context, account model.AccountBase) (model.AccountBase, error) {
	row, err := r.sqlc.CreateAccount(ctx, sqlc.CreateAccountParams{
		Username: account.Username,
		Email:    account.Email,
		Name:     account.Name,
		Password: account.Password,
	})
	if err != nil {
		return model.AccountBase{}, err
	}

	return model.AccountBase{
		ID:        row.ID,
		Username:  row.Username,
		Email:     row.Email,
		Name:      row.Name,
		Password:  row.Password,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

type UpdateAccountParams struct {
	ID       int64
	Username *string
	Email    *string
	Name     *string
	Password *string
}

func (r *RepositoryImpl) UpdateAccount(ctx context.Context, params UpdateAccountParams) (model.AccountBase, error) {
	row, err := r.sqlc.UpdateAccount(ctx, sqlc.UpdateAccountParams{
		ID:       params.ID,
		Username: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Username),
		Email:    *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Email),
		Name:     *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name),
		Password: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Password),
	})
	if err != nil {
		return model.AccountBase{}, err
	}

	return model.AccountBase{
		ID:        row.ID,
		Username:  row.Username,
		Name:      row.Name,
		Email:     row.Email,
		Password:  row.Password,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (r *RepositoryImpl) DeleteAccount(ctx context.Context, accountID int64) error {
	return r.sqlc.DeleteAccount(ctx, accountID)
}
