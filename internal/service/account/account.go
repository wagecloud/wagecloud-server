package account

import (
	"context"
	"fmt"

	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/repository"
	"github.com/wagecloud/wagecloud-server/internal/util/jwt"
	"golang.org/x/crypto/bcrypt"
)

var _ ServiceInterface = (*Service)(nil)

type Service struct {
	repo repository.Repository
}

type ServiceInterface interface {
	GetAccount(ctx context.Context, params GetAccountParams) (model.AccountUser, error)
	LoginUser(ctx context.Context, params LoginUserParams) (LoginUserResult, error)
	RegisterUser(ctx context.Context, account model.AccountUser) (RegisterUserResult, error)
}

func NewService(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

type GetAccountParams struct {
	ID       *int64
	Username *string
	Email    *string
}

func (s *Service) GetAccount(ctx context.Context, params GetAccountParams) (model.AccountUser, error) {
	account, err := s.repo.GetAccount(ctx, repository.GetAccountParams{
		ID:       params.ID,
		Username: params.Username,
		Email:    params.Email,
	})
	if err != nil {
		return model.AccountUser{}, err
	}

	return model.AccountUser{
		AccountBase: account,
	}, nil
}

type LoginUserParams struct {
	ID       *int64
	Username *string
	Email    *string
	Password string
}

type LoginUserResult struct {
	Token   string
	Account model.AccountUser
}

func (s *Service) LoginUser(ctx context.Context, params LoginUserParams) (LoginUserResult, error) {
	account, err := s.repo.GetAccount(ctx, repository.GetAccountParams{
		ID:       params.ID,
		Username: params.Username,
		Email:    params.Email,
	})
	if err != nil {
		return LoginUserResult{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(params.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return LoginUserResult{}, fmt.Errorf("wrong password")
		}
		return LoginUserResult{}, fmt.Errorf("failed to compare password: %w", err)
	}

	token, err := jwt.GenerateAccessToken(account)
	if err != nil {
		return LoginUserResult{}, err
	}

	return LoginUserResult{
		Token: token,
		Account: model.AccountUser{
			AccountBase: account,
		},
	}, nil
}

type RegisterUserResult struct {
	Token   string
	Account model.AccountUser
}

func (s *Service) RegisterUser(ctx context.Context, account model.AccountUser) (res RegisterUserResult, err error) {
	txRepo, err := s.repo.Begin(ctx)
	if err != nil {
		return res, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer txRepo.Rollback(ctx)

	// Role must set to USER
	account.Role = model.RoleUser

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return res, fmt.Errorf("failed to hash password: %w", err)
	}
	// Set the hashed password in the account
	account.Password = string(hashedPassword)

	createdAccount, err := txRepo.CreateAccount(ctx, account.Base())
	if err != nil {
		return res, fmt.Errorf("failed to create account: %w", err)
	}

	token, err := jwt.GenerateAccessToken(createdAccount)
	if err != nil {
		return res, fmt.Errorf("failed to generate access token: %w", err)
	}

	if err = txRepo.Commit(ctx); err != nil {
		return res, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return RegisterUserResult{
		Token: token,
		Account: model.AccountUser{
			AccountBase: createdAccount,
		},
	}, nil
}
