package accountsvc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wagecloud/wagecloud-server/config"
	accountmodel "github.com/wagecloud/wagecloud-server/internal/modules/account/model"
	accountstorage "github.com/wagecloud/wagecloud-server/internal/modules/account/storage"
	"github.com/wagecloud/wagecloud-server/internal/utils/cache"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenHeader = "authorization"
	tokenPrefix = "Bearer "
)

var (
	// TODO: use redis
	claimsCache = cache.NewCache[string, accountmodel.Claims]()
)

type ServiceImpl struct {
	storage *accountstorage.Storage
}

type Service interface {
	GetAccount(ctx context.Context, params GetAccountParams) (accountmodel.AccountUser, error)
	LoginUser(ctx context.Context, params LoginUserParams) (LoginUserResult, error)
	RegisterUser(ctx context.Context, account accountmodel.AccountUser) (RegisterUserResult, error)
}

func NewService(storage *accountstorage.Storage) *ServiceImpl {
	return &ServiceImpl{
		storage: storage,
	}
}

type GetAccountParams struct {
	ID       *int64
	Username *string
	Email    *string
}

func (s *ServiceImpl) GetAccount(ctx context.Context, params GetAccountParams) (accountmodel.AccountUser, error) {
	account, err := s.storage.GetAccount(ctx, accountstorage.GetAccountParams{
		ID:       params.ID,
		Username: params.Username,
		Email:    params.Email,
	})
	if err != nil {
		return accountmodel.AccountUser{}, err
	}

	return accountmodel.AccountUser{
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
	Account accountmodel.AccountUser
}

func (s *ServiceImpl) LoginUser(ctx context.Context, params LoginUserParams) (LoginUserResult, error) {
	account, err := s.storage.GetAccount(ctx, accountstorage.GetAccountParams{
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

	token, err := GenerateAccessToken(account)
	if err != nil {
		return LoginUserResult{}, err
	}

	return LoginUserResult{
		Token: token,
		Account: accountmodel.AccountUser{
			AccountBase: account,
		},
	}, nil
}

type RegisterUserResult struct {
	Token   string
	Account accountmodel.AccountUser
}

func (s *ServiceImpl) RegisterUser(ctx context.Context, account accountmodel.AccountUser) (res RegisterUserResult, err error) {
	txStorage, err := s.storage.BeginTx(ctx)
	if err != nil {
		return res, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer txStorage.Rollback(ctx)

	// Role must set to USER
	account.Role = accountmodel.RoleUser

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return res, fmt.Errorf("failed to hash password: %w", err)
	}
	// Set the hashed password in the account
	account.Password = string(hashedPassword)

	createdAccount, err := txStorage.CreateAccount(ctx, account.Base())
	if err != nil {
		return res, fmt.Errorf("failed to create account: %w", err)
	}

	createdUser, err := txStorage.CreateUser(ctx, accountmodel.AccountUser{
		AccountBase: createdAccount,
		Email:       account.Email,
	})
	if err != nil {
		return res, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := GenerateAccessToken(createdAccount)
	if err != nil {
		return res, fmt.Errorf("failed to generate access token: %w", err)
	}

	if err = txStorage.Commit(ctx); err != nil {
		return res, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return RegisterUserResult{
		Token: token,
		Account: accountmodel.AccountUser{
			AccountBase: createdAccount,
			Email:       createdUser.Email,
		},
	}, nil
}

func GenerateAccessToken(account accountmodel.Account) (string, error) {
	tokenDuration := time.Duration(config.GetConfig().App.AccessTokenDuration * int64(time.Second))

	claims := accountmodel.Claims{
		AccountID: account.Base().ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "wagecloud",
			Subject:   strconv.Itoa(int(account.Base().ID)),
			Audience:  []string{"wagecloud"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := config.GetConfig().SensitiveKeys.JWTSecret

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	// test if token is valid
	_, err = ValidateAccessToken(signedToken)
	if err != nil {
		return "", fmt.Errorf("failed to validate token: %w", err)
	}

	return signedToken, nil
}

func ValidateAccessToken(tokenStr string) (claims accountmodel.Claims, err error) {
	secret := config.GetConfig().SensitiveKeys.JWTSecret

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return claims, err
	}

	if !token.Valid {
		return claims, errors.New("invalid token or token expired")
	}

	return claims, nil
}

func GetClaims(r *http.Request) (claims accountmodel.Claims, err error) {
	token := r.Header.Get(tokenHeader)

	claims, ok := claimsCache.Get(token)
	if ok {
		return claims, nil
	}

	claims, err = ValidateAccessToken(strings.TrimPrefix(token, tokenPrefix))
	if err != nil {
		return accountmodel.Claims{}, err
	}

	claimsCache.Set(token, claims, 5*60*time.Second)
	return claims, nil
}
