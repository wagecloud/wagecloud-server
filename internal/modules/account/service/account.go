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
	"github.com/patrickmn/go-cache"
	"github.com/wagecloud/wagecloud-server/config"
	accountmodel "github.com/wagecloud/wagecloud-server/internal/modules/account/model"
	accountstorage "github.com/wagecloud/wagecloud-server/internal/modules/account/storage"
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
}

func (s *ServiceImpl) GetUser(ctx context.Context, params GetUserParams) (accountmodel.AccountUser, error) {
	user, err := s.storage.GetUser(ctx, accountstorage.GetUserParams{
		ID:       params.ID,
		Username: params.Username,
		Email:    params.Email,
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
	Password string
}

type LoginUserResult struct {
	Token   string
	Account accountmodel.AccountUser
}

func (s *ServiceImpl) LoginUser(ctx context.Context, params LoginUserParams) (LoginUserResult, error) {
	user, err := s.storage.GetUser(ctx, accountstorage.GetUserParams{
		ID:       params.ID,
		Username: params.Username,
		Email:    params.Email,
	})
	if err != nil {
		return LoginUserResult{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return LoginUserResult{}, fmt.Errorf("wrong password")
		}
		return LoginUserResult{}, fmt.Errorf("failed to compare password: %w", err)
	}

	token, err := GenerateAccessToken(user)
	if err != nil {
		return LoginUserResult{}, err
	}

	return LoginUserResult{
		Token:   token,
		Account: user,
	}, nil
}

type RegisterUserParams struct {
	Name     string
	Email    string
	Username string
	Password string
}

type RegisterUserResult struct {
	Token   string
	Account accountmodel.AccountUser
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

	createdAccount, err := txStorage.CreateAccount(ctx, accountmodel.AccountBase{
		Type:     accountmodel.AccountTypeUser,
		Name:     params.Name,
		Username: params.Username,
		Password: string(hashedPassword),
	})
	if err != nil {
		return res, fmt.Errorf("failed to create account: %w", err)
	}

	createdUser, err := txStorage.CreateUser(ctx, accountmodel.AccountUser{
		AccountBase: createdAccount,
		Email:       params.Email,
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
