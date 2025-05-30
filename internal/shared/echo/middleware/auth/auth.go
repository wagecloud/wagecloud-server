package auth

import (
	"net/http"
	"strings"
	"time"

	accountmodel "github.com/wagecloud/wagecloud-server/internal/modules/account/model"
	accountsvc "github.com/wagecloud/wagecloud-server/internal/modules/account/service"
	"github.com/wagecloud/wagecloud-server/internal/utils/cache"
)

const (
	tokenHeader = "authorization"
	tokenPrefix = "Bearer "
)

var (
	claimsCache = cache.NewCache[string, accountmodel.Claims]()
)

func GetClaims(r *http.Request) (claims accountmodel.Claims, err error) {
	token := r.Header.Get(tokenHeader)

	claims, ok := claimsCache.Get(token)
	if ok {
		return claims, nil
	}

	claims, err = accountsvc.ValidateAccessToken(strings.TrimPrefix(token, tokenPrefix))
	if err != nil {
		return accountmodel.Claims{}, err
	}

	claimsCache.Set(token, claims, 5*60*time.Second)
	return claims, nil
}
