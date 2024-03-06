package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (ts *PostgresTestSuite) TestLogin() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	bUsername := make([]byte, 10)
	_, err := rand.Read(bUsername)
	ts.Require().NoError(err)

	bPassword := make([]byte, 10)
	_, err = rand.Read(bPassword)
	ts.Require().NoError(err)

	data := models.User{
		Username: hex.EncodeToString(bUsername),
		Password: hex.EncodeToString(bPassword),
	}
	id, err := ts.storage.Add(ctx, data)
	ts.Require().NoError(err)

	userID, err := ts.storage.Login(ctx, data)
	ts.NoError(err)
	ts.EqualValues(id, userID)
}
