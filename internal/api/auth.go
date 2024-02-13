package api

import (
	"net/http"

	uService "github.com/fishus/go-advanced-gophermart/internal/service/user"
)

func (s *server) auth(r *http.Request) (*uService.JWTClaims, error) {
	auth := r.Header.Get("Authorization")
	return s.service.User().CheckAuthorizationHeader(auth)
}
