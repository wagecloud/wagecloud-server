package accountstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wagecloud/wagecloud-server/gen/sqlc"
	"github.com/wagecloud/wagecloud-server/internal/client/pgxpool"
	accountmodel "github.com/wagecloud/wagecloud-server/internal/modules/account/model"
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
	// Add more tx fields if needed (e.g., redis, mongo, etc.)
}

func NewStorage(db pgxpool.DBTX) *Storage {
	return &Storage{
		db:   db,
		sqlc: sqlc.New(db),
	}
}

// BeginTx starts a pseudo nested transaction.
func (s *Storage) BeginTx(ctx context.Context) (*TxStorage, error) {
	// Add more tx begin logics if needed (e.g., redis, mongo, etc.)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &TxStorage{
		Storage: NewStorage(tx),
		tx:      tx,
	}, nil
}

// Commit commits the transaction if this is a real transaction or releases the savepoint if this is a pseudo nested
// transaction. Commit will return an error if the Tx is already closed, but is otherwise safe to call multiple times.
func (ts *TxStorage) Commit(ctx context.Context) error {
	// Add more tx commit logics if needed (e.g., redis, mongo, etc.)
	return ts.tx.Commit(ctx)
}

// Rollback rolls back the transaction. Rollback will return an error if the Tx is already closed, but is otherwise safe
// to call multiple times. Hence, a defer storage.Rollback() is safe (must safe) even if storage.Commit() will be
// called first in a non-error condition. Any other failure of a real transaction will result in the connection being closed.
func (ts *TxStorage) Rollback(ctx context.Context) error {
	// Add more tx rollback logics if needed (e.g., redis, mongo, etc.)
	return ts.tx.Rollback(ctx)
}

type GetAccountParams struct {
	ID       *int64
	Username *string
	Email    *string
}

func (s *Storage) GetAccount(ctx context.Context, params GetAccountParams) (accountmodel.AccountBase, error) {
	if params.ID == nil && params.Username == nil {
		return accountmodel.AccountBase{}, fmt.Errorf("at least one of ID, Username must be provided")
	}

	row, err := s.sqlc.GetAccount(ctx, sqlc.GetAccountParams{
		ID:       *pgxptr.PtrToPgtype(&pgtype.Int8{}, params.ID),
		Username: *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Username),
		Email:    *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Email),
	})
	if err != nil {
		return accountmodel.AccountBase{}, err
	}

	return accountmodel.AccountBase{
		ID:        row.ID,
		Type:      accountmodel.AccountType(row.Type),
		Name:      row.Name,
		Username:  row.Username,
		Password:  row.Password,
		CreatedAt: row.CreatedAt.Time.UnixMilli(),
		UpdatedAt: row.UpdatedAt.Time.UnixMilli(),
	}, nil
}

type ListAccountsParams struct {
	pagination.PaginationParams
	ID            *string
	Type          *accountmodel.AccountType
	Username      *string
	Name          *string
	CreatedAtFrom *int64
	CreatedAtTo   *int64
}

func (s *Storage) CountAccounts(ctx context.Context, params ListAccountsParams) (int64, error) {
	return s.sqlc.CountAccounts(ctx, sqlc.CountAccountsParams{
		ID:            *pgxptr.PtrToPgtype(&pgtype.Text{}, params.ID),
		Type:          *pgxptr.PtrBrandedToPgType(&sqlc.NullAccountType{}, params.Type),
		Username:      *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Username),
		Name:          *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Name),
		CreatedAtFrom: *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtFrom)),
		CreatedAtTo:   *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtTo)),
	})
}

func (s *Storage) ListAccounts(ctx context.Context, params ListAccountsParams) ([]accountmodel.AccountBase, error) {
	rows, err := s.sqlc.ListAccounts(ctx, sqlc.ListAccountsParams{
		Limit:         params.Limit,
		Offset:        params.Offset(),
		ID:            *pgxptr.PtrToPgtype(&pgtype.Text{}, params.ID),
		Type:          *pgxptr.PtrBrandedToPgType(&sqlc.NullAccountType{}, params.Type),
		Name:          *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Name),
		Username:      *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Username),
		CreatedAtFrom: *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtFrom)),
		CreatedAtTo:   *pgxptr.PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.CreatedAtTo)),
	})
	if err != nil {
		return nil, err
	}

	var accounts []accountmodel.AccountBase
	for _, row := range rows {
		accounts = append(accounts, accountmodel.AccountBase{
			ID:        row.ID,
			Type:      accountmodel.AccountType(row.Type),
			Name:      row.Name,
			Username:  row.Username,
			Password:  row.Password,
			CreatedAt: row.CreatedAt.Time.UnixMilli(),
			UpdatedAt: row.UpdatedAt.Time.UnixMilli(),
		})
	}

	return accounts, nil
}

func (s *Storage) CreateAccount(ctx context.Context, account accountmodel.AccountBase) (accountmodel.AccountBase, error) {
	row, err := s.sqlc.CreateAccount(ctx, sqlc.CreateAccountParams{
		Type:     sqlc.AccountType(account.Type),
		Name:     account.Name,
		Username: account.Username,
		Password: account.Password,
	})
	if err != nil {
		return accountmodel.AccountBase{}, err
	}

	return accountmodel.AccountBase{
		ID:        row.ID,
		Type:      accountmodel.AccountType(row.Type),
		Name:      row.Name,
		Username:  row.Username,
		Password:  row.Password,
		CreatedAt: row.CreatedAt.Time.UnixMilli(),
		UpdatedAt: row.UpdatedAt.Time.UnixMilli(),
	}, nil
}

type UpdateAccountParams struct {
	ID       int64
	Username *string
	Name     *string
	Password *string
}

func (s *Storage) UpdateAccount(ctx context.Context, params UpdateAccountParams) (accountmodel.AccountBase, error) {
	row, err := s.sqlc.UpdateAccount(ctx, sqlc.UpdateAccountParams{
		ID:       params.ID,
		Name:     *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Name),
		Username: *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Username),
		Password: *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Password),
	})
	if err != nil {
		return accountmodel.AccountBase{}, err
	}

	return accountmodel.AccountBase{
		ID:        row.ID,
		Type:      accountmodel.AccountType(row.Type),
		Name:      row.Name,
		Username:  row.Username,
		Password:  row.Password,
		CreatedAt: row.CreatedAt.Time.UnixMilli(),
		UpdatedAt: row.UpdatedAt.Time.UnixMilli(),
	}, nil
}

func (s *Storage) DeleteAccount(ctx context.Context, accountID int64) error {
	return s.sqlc.DeleteAccount(ctx, accountID)
}

type GetUserParams struct {
	ID       *int64
	Username *string
	Email    *string
}

func (s *Storage) GetUser(ctx context.Context, params GetUserParams) (accountmodel.AccountUser, error) {
	if params.ID == nil && params.Username == nil && params.Email == nil {
		return accountmodel.AccountUser{}, fmt.Errorf("at least one of ID, Username, Email must be provided")
	}

	row, err := s.sqlc.GetUser(ctx, sqlc.GetUserParams{
		ID:       *pgxptr.PtrToPgtype(&pgtype.Int8{}, params.ID),
		Username: *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Username),
		Email:    *pgxptr.PtrToPgtype(&pgtype.Text{}, params.Email),
	})
	if err != nil {
		return accountmodel.AccountUser{}, err
	}

	return accountmodel.AccountUser{
		AccountBase: accountmodel.AccountBase{
			ID:        row.ID,
			Type:      accountmodel.AccountType(row.Type),
			Name:      row.Name,
			Username:  row.Username,
			Password:  row.Password,
			CreatedAt: row.CreatedAt.Time.UnixMilli(),
			UpdatedAt: row.UpdatedAt.Time.UnixMilli(),
		},
		Email: row.Email,
	}, nil
}

func (s *Storage) CreateUser(ctx context.Context, user accountmodel.AccountUser) (accountmodel.AccountUser, error) {
	row, err := s.sqlc.CreateUser(ctx, sqlc.CreateUserParams{
		ID:    user.ID,
		Email: user.Email,
	})
	if err != nil {
		return accountmodel.AccountUser{}, err
	}

	return accountmodel.AccountUser{
		AccountBase: accountmodel.AccountBase{
			ID:   row.ID,
			Type: accountmodel.AccountTypeUser,
		},
		Email: row.Email,
	}, nil
}
