package user

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

var ErrInvalidToken = errors.New("invalid token")

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

func (s *service) DecryptToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return []byte(s.cfg.JWTSecretKey), nil
		})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *service) CheckAuthorizationHeader(auth string) (*JWTClaims, error) {
	if len(auth) < 7 {
		return nil, errors.New("invalid bearer authorization header")
	}

	authType := strings.ToLower(auth[:6])
	if authType != "bearer" {
		return nil, errors.New("invalid bearer authorization header")
	}

	token, err := s.DecryptToken(auth[7:])
	if err != nil {
		return nil, ErrInvalidToken
	}

	return token, nil
}
