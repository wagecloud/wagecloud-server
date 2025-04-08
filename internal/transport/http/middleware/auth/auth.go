package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/util/cache"
	"github.com/wagecloud/wagecloud-server/internal/util/jwt"
)

const (
	tokenHeader = "authorization"
	tokenPrefix = "Bearer "
)

var (
	claimsCache = cache.NewCache[string, model.Claims]()
)

// func NewAuthMiddleware() func(next http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			GetClaims(r) // Invoke just for

// 			// TODO: làm config cho cache time
// 			claimsCache.Set(r.Header.Get(tokenHeader), claims, 5*60*time.Second)

// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

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
