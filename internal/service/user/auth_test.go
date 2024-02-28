package user

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (ts *UserServiceTestSuite) TestBuildToken() {
	userID := models.UserID(uuid.New().String())
	tokenString, err := ts.service.BuildToken(userID)
	ts.Require().NoError(err)

	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(ts.cfg.JWTSecretKey), nil
	})
	ts.NoError(err)
	ts.Equal(token.Valid, true)
	ts.Equal(claims.UserID, userID)
}

func (ts *UserServiceTestSuite) TestDecryptToken() {
	userID := models.UserID(uuid.New().String())
	tokenString, err := ts.service.BuildToken(userID)
	ts.Require().NoError(err)

	claims, err := ts.service.DecryptToken(tokenString)
	ts.NoError(err)
	ts.Equal(claims.UserID, userID)
}

func (ts *UserServiceTestSuite) TestCheckAuthorizationHeader() {
	userID := models.UserID(uuid.New().String())
	tokenString, err := ts.service.BuildToken(userID)
	ts.Require().NoError(err)

	auth := "Bearer " + tokenString

	claims, err := ts.service.CheckAuthorizationHeader(auth)
	ts.NoError(err)
	ts.Equal(claims.UserID, userID)
}
