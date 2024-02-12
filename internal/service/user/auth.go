package user

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID models.UserID
}

// BuildToken build token for user
func (s *service) BuildToken(userID models.UserID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.JWTExpires)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
