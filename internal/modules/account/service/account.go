package accountsvc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/patrickmn/go-cache"
	"github.com/wagecloud/wagecloud-server/config"
	accountmodel "github.com/wagecloud/wagecloud-server/internal/modules/account/model"
	accountstorage "github.com/wagecloud/wagecloud-server/internal/modules/account/storage"
	"github.com/wagecloud/wagecloud-server/internal/utils/ptr"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenHeader        = "authorization"
	tokenPrefix        = "Bearer "
	tokenCacheDuration = 5 * 60 * time.Second
)

var (
	claimsCache = cache.New(tokenCacheDuration, 10*time.Minute) // 5 min default expiration, 10 min cleanup interval
)

type ServiceImpl struct {
	storage *accountstorage.Storage
}

type Service interface {
	UpdateAccount(ctx context.Context, params UpdateAccountParams) (accountmodel.AccountBase, error)
	UpdateUser(ctx context.Context, params UpdateUserParams) (accountmodel.AccountUser, error)
	GetUser(ctx context.Context, params GetUserParams) (accountmodel.AccountUser, error)
	LoginUser(ctx context.Context, params LoginUserParams) (LoginUserResult, error)
	RegisterUser(ctx context.Context, params RegisterUserParams) (RegisterUserResult, error)
}

func NewService(storage *accountstorage.Storage) Service {
	return &ServiceImpl{
		storage: storage,
	}
}

type GetUserParams struct {
	Account  accountmodel.AuthenticatedAccount
	ID       *int64
	Username *string
	Email    *string
	Phone    *string
}

type UpdateAccountParams struct {
	Account         accountmodel.AuthenticatedAccount
	CurrentPassword string
	Username        *string
	NewPassword     *string
}

func (s *ServiceImpl) UpdateAccount(ctx context.Context, params UpdateAccountParams) (accountmodel.AccountBase, error) {
	account, err := s.storage.GetAccount(ctx, accountstorage.GetAccountParams{
		Type: accountmodel.AccountTypeUser,
		ID:   &params.Account.AccountID,
	})
	if err != nil {
		// TODO: move this sql handle error to storage layer, because sql come from there
		if errors.Is(err, sql.ErrNoRows) {
			return accountmodel.AccountBase{}, accountmodel.ErrAccountNotFound
		}
		return accountmodel.AccountBase{}, fmt.Errorf("failed to get account: %w", err)
	}

	// Check current password
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(params.CurrentPassword)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return accountmodel.AccountBase{}, accountmodel.ErrWrongCurrentPassword
		}
		return accountmodel.AccountBase{}, fmt.Errorf("failed to compare current password: %w", err)
	}

	if params.NewPassword != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*params.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return accountmodel.AccountBase{}, fmt.Errorf("failed to hash password: %w", err)
		}

		fmt.Println("New password %s", *params.NewPassword)

		params.NewPassword = ptr.ToPtr(string(hashedPassword))

		fmt.Println("Updating account with new password")
		fmt.Printf("New password hash: %s\n", *params.NewPassword)

	}

	updatedAccount, err := s.storage.UpdateAccount(ctx, accountstorage.UpdateAccountParams{
		ID:       params.Account.AccountID,
		Username: params.Username,
		Password: params.NewPassword,
	})
	if err != nil {
		return accountmodel.AccountBase{}, err
	}

	return updatedAccount, nil
}

type UpdateUserParams struct {
	ID          int64
	FirstName   *string
	LastName    *string
	Email       *string
	NullEmail   bool
	Phone       *string
	NullPhone   bool
	Company     *string
	NullCompany bool
	Address     *string
	NullAddress bool
}

func (s *ServiceImpl) UpdateUser(ctx context.Context, params UpdateUserParams) (accountmodel.AccountUser, error) {
	updatedUser, err := s.storage.UpdateUser(ctx, accountstorage.UpdateUserParams{
		ID:          params.ID,
		FirstName:   params.FirstName,
		LastName:    params.LastName,
		Email:       params.Email,
		NullEmail:   params.NullEmail,
		Phone:       params.Phone,
		NullPhone:   params.NullPhone,
		Company:     params.Company,
		NullCompany: params.NullCompany,
		Address:     params.Address,
		NullAddress: params.NullAddress,
	})
	if err != nil {
		return accountmodel.AccountUser{}, err
	}

	return updatedUser, nil
}

func (s *ServiceImpl) GetUser(ctx context.Context, params GetUserParams) (accountmodel.AccountUser, error) {
	user, err := s.storage.GetUser(ctx, accountstorage.GetUserParams{
		ID:       params.ID,
		Username: params.Username,
		Email:    params.Email,
		Phone:    params.Phone,
	})
	if err != nil {
		return accountmodel.AccountUser{}, err
	}

	if err := s.canAccess(ctx, canAccessParams{
		Account:   params.Account,
		AccountID: user.ID,
	}); err != nil {
		return accountmodel.AccountUser{}, err
	}

	return user, nil
}

type LoginUserParams struct {
	ID       *int64
	Username *string
	Email    *string
	Phone    *string
	Password string
}

type LoginUserResult struct {
	Token   string                   `json:"token"`
	Account accountmodel.AccountBase `json:"account"`
}

func (s *ServiceImpl) LoginUser(ctx context.Context, params LoginUserParams) (LoginUserResult, error) {
	account, err := s.storage.GetAccount(ctx, accountstorage.GetAccountParams{
		Type:     accountmodel.AccountTypeUser,
		ID:       params.ID,
		Username: params.Username,
		Email:    params.Email,
		Phone:    params.Phone,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LoginUserResult{}, accountmodel.ErrInvalidCredentials
		}
		return LoginUserResult{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(params.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return LoginUserResult{}, accountmodel.ErrInvalidCredentials
		}
		return LoginUserResult{}, fmt.Errorf("failed to compare password: %w", err)
	}

	token, err := GenerateAccessToken(account.ID)
	if err != nil {
		return LoginUserResult{}, err
	}

	return LoginUserResult{
		Token:   token,
		Account: account,
	}, nil
}

type RegisterUserParams struct {
	FirstName string
	LastName  string
	Username  string
	Password  string
	Email     *string
	Phone     *string
}

type RegisterUserResult struct {
	Token   string                   `json:"token"`
	Account accountmodel.AccountUser `json:"account"`
}

func (s *ServiceImpl) RegisterUser(ctx context.Context, params RegisterUserParams) (res RegisterUserResult, err error) {
	txStorage, err := s.storage.BeginTx(ctx)
	if err != nil {
		return res, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer txStorage.Rollback(ctx)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return res, fmt.Errorf("failed to hash password: %w", err)
	}

	if _, err := txStorage.GetAccount(ctx, accountstorage.GetAccountParams{
		Type:     accountmodel.AccountTypeUser,
		Username: &params.Username,
		Email:    params.Email,
		Phone:    params.Phone,
	}); err == nil {
		return res, accountmodel.ErrAccountAlreadyExists
	} else {
		if !errors.Is(err, sql.ErrNoRows) {
			return res, fmt.Errorf("failed to check existing account: %w", err)
		}
	}

	createdAccount, err := txStorage.CreateAccount(ctx, accountmodel.AccountBase{
		Type:     accountmodel.AccountTypeUser,
		Username: params.Username,
		Password: string(hashedPassword),
	})
	if err != nil {
		return res, fmt.Errorf("failed to create account: %w", err)
	}

	createdUser, err := txStorage.CreateUser(ctx, accountmodel.AccountUser{
		ID:        createdAccount.ID,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Phone:     params.Phone,
	})
	if err != nil {
		return res, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := GenerateAccessToken(createdAccount.ID)
	if err != nil {
		return res, fmt.Errorf("failed to generate access token: %w", err)
	}

	if err = txStorage.Commit(ctx); err != nil {
		return res, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return RegisterUserResult{
		Token: token,
		Account: accountmodel.AccountUser{
			FirstName: createdUser.FirstName,
			LastName:  createdUser.LastName,
			Email:     createdUser.Email,
			Phone:     createdUser.Phone,
		},
	}, nil
}

func GenerateAccessToken(accountID int64) (string, error) {
	tokenDuration := time.Duration(config.GetConfig().App.AccessTokenDuration * int64(time.Second))

	claims := accountmodel.Claims{
		AccountID: accountID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "wagecloud",
			Subject:   strconv.Itoa(int(accountID)),
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

// GetClaims retrieves and validates JWT claims from the token, using an in-memory cache
func GetClaims(r *http.Request) (claims accountmodel.Claims, err error) {
	token := r.Header.Get(tokenHeader)

	if token == "" {
		return accountmodel.Claims{}, fmt.Errorf("missing authorization header")
	}

	// Try to get claims from cache first
	if cachedClaims, found := claimsCache.Get(token); found {
		if claims, ok := cachedClaims.(accountmodel.Claims); ok {
			return claims, nil
		}
	}

	// If not in cache, validate token and store in cache
	claims, err = ValidateAccessToken(strings.TrimPrefix(token, tokenPrefix))
	if err != nil {
		return accountmodel.Claims{}, err
	}

	// Store claims in cache
	claimsCache.Set(token, claims, tokenCacheDuration)

	return claims, nil
}

type canAccessParams struct {
	Account   accountmodel.AuthenticatedAccount
	AccountID int64
}

func (s *ServiceImpl) canAccess(_ context.Context, params canAccessParams) error {
	if params.Account.Type == accountmodel.AccountTypeAdmin {
		return nil
	}

	if params.Account.AccountID == params.AccountID {
		return nil
	}

	return fmt.Errorf("access denied: account %d cannot access account %d", params.Account.AccountID, params.AccountID)
}
