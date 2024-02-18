package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (ts *PostgresTestSuite) TestUserByID() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	ts.Run("Return user by ID", func() {
		bUsername := make([]byte, 10)
		_, err := rand.Read(bUsername)
		ts.Require().NoError(err)

		bPassword := make([]byte, 10)
		_, err = rand.Read(bPassword)
		ts.Require().NoError(err)

		data := models.User{
			Username:  hex.EncodeToString(bUsername),
			Password:  hex.EncodeToString(bPassword),
			CreatedAt: time.Now().UTC().Round(5 * time.Second),
		}
		id, err := ts.storage.UserAdd(ctx, data)
		ts.Require().NoError(err)
		data.ID = id
		data.Password = "" // password is always empty

		user, err := ts.storage.UserByID(ctx, data.ID)
		ts.NoError(err)
		user.CreatedAt = user.CreatedAt.UTC().Round(5 * time.Second)
		ts.EqualValues(data, user)
	})

	ts.Run("User not found", func() {
		userID := models.UserID(uuid.New().String())
		_, err := ts.storage.UserByID(ctx, userID)
		ts.Error(err)
		ts.ErrorIs(err, store.ErrNotFound)
	})
}
