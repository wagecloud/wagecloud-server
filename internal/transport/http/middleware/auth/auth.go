package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/util/cache"
	"github.com/wagecloud/wagecloud-server/internal/util/jwt"
)

// type ctxKey string

const (
	tokenHeader = "authorization"
	tokenPrefix = "Bearer "
	// CtxClaims   ctxKey = "ctx-claims" // Storing model.Claims in context
	// CtxToken    ctxKey = "ctx-token"  // Storing token in context
)

var (
	claimsCache = cache.NewCache[string, model.Claims]()
)

func GetClaims(r *http.Request) (claims model.Claims, err error) {
	token := r.Header.Get(tokenHeader)

	claims, ok := claimsCache.Get(token)
	if ok {
		return claims, nil
	}

	claims, err = jwt.ValidateAccessToken(strings.TrimPrefix(token, tokenPrefix))
	if err != nil {
		return model.Claims{}, err
	}

	claimsCache.Set(token, claims, 5*60*time.Second)
	return claims, nil
}
