package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (ts *PostgresTestSuite) TestUserAdd() {
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

	ts.Run("Positive case", func() {
		userID, err := ts.storage.UserAdd(ctx, data)
		ts.NoError(err)

		var want struct {
			id        string
			username  string
			password  string
			createdAt time.Time
		}
		err = ts.pool.QueryRow(ctx, "SELECT id, username, password, created_at FROM users WHERE id = @id;",
			pgx.NamedArgs{"id": userID}).Scan(&want.id, &want.username, &want.password, &want.createdAt)
		ts.NoError(err)
		ts.Equal(userID.String(), want.id)
		ts.Equal(data.Username, want.username)
		ts.NotEmpty(want.password)
		ts.Equal(time.Now().UTC().Round(10*time.Second), want.createdAt.Round(10*time.Second))
	})

	ts.Run("DuplicateUser", func() {
		_, err = ts.storage.UserAdd(ctx, data)
		ts.Error(err)
		ts.ErrorIs(err, store.ErrAlreadyExists)
	})

	ts.Run("IncorrectUser", func() {
		data := &models.User{
			Username: "",
			Password: "",
		}
		_, err := ts.storage.UserAdd(ctx, *data)
		ts.Error(err)
		ts.ErrorIs(err, store.ErrIncorrectData)
	})
}
