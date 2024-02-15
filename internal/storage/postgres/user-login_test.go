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

func (ts *PostgresTestSuite) TestUserLogin() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	bUsername := make([]byte, 10)
	_, err := rand.Read(bUsername)
	ts.Require().NoError(err)

	bPassword := make([]byte, 10)
	_, err = rand.Read(bPassword)
	ts.Require().NoError(err)

	id := models.UserID(uuid.New().String())
	data := &models.User{
		Username:  hex.EncodeToString(bUsername),
		Password:  hex.EncodeToString(bPassword),
		CreatedAt: time.Now().UTC().Round(1 * time.Second),
	}

	_, err = ts.pool.Exec(ctx, `INSERT INTO users (id, username, password, created_at) VALUES (@id, @username, crypt(@password, gen_salt('bf')), @created_at);`,
		pgx.NamedArgs{
			"id":         id,
			"username":   data.Username,
			"password":   data.Password,
			"created_at": data.CreatedAt,
		})
	ts.NoError(err)

	userID, err := ts.storage.UserLogin(ctx, *data)
	ts.NoError(err)
	ts.EqualValues(id, userID)
}
