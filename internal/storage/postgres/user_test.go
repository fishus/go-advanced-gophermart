package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

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
		ID:        models.UserID(uuid.New().String()),
		Username:  hex.EncodeToString(bUsername),
		Password:  hex.EncodeToString(bPassword),
		CreatedAt: time.Now().UTC().Round(1 * time.Second),
	}

	_, err = ts.pool.Exec(ctx, `INSERT INTO users (id, username, password, created_at) VALUES (@id, @username, @password, @createdAt);`,
		pgx.NamedArgs{
			"id":        data.ID,
			"username":  data.Username,
			"password":  data.Password,
			"createdAt": data.CreatedAt,
		})
	ts.NoError(err)

	data.Password = "" // password is always empty
	user, err := ts.storage.UserByID(ctx, data.ID)
	ts.NoError(err)
	ts.EqualValues(data, user)
}
