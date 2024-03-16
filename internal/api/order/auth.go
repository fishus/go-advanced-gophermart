package order

import (
	"net/http"

	uService "github.com/fishus/go-advanced-gophermart/internal/service/user"
)

func (a *api) auth(r *http.Request) (*uService.JWTClaims, error) {
	auth := r.Header.Get("Authorization")
	return a.service.User().CheckAuthorizationHeader(auth)
}
