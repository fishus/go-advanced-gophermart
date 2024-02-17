package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (ts *PostgresTestSuite) TestUserByID() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	bUsername := make([]byte, 10)
	_, err := rand.Read(bUsername)
	ts.Require().NoError(err)

	bPassword := make([]byte, 10)
	_, err = rand.Read(bPassword)
	ts.Require().NoError(err)

	data := models.User{
		Username:  hex.EncodeToString(bUsername),
		Password:  hex.EncodeToString(bPassword),
		CreatedAt: time.Now().UTC().Round(1 * time.Second),
	}
	id, err := ts.storage.UserAdd(ctx, data)
	ts.Require().NoError(err)
	data.ID = id
	data.Password = "" // password is always empty

	user, err := ts.storage.UserByID(ctx, data.ID)
	ts.NoError(err)
	user.CreatedAt = user.CreatedAt.UTC().Round(1 * time.Second)
	ts.EqualValues(data, user)
}
