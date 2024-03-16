package user

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (s *storage) Add(ctx context.Context, user models.User) (userID models.UserID, err error) {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	if user.Username == "" || user.Password == "" {
		err = store.ErrIncorrectData
		return
	}

	err = s.pool.QueryRow(ctxQuery, `INSERT INTO users (username, password) VALUES (@username, crypt(@password, gen_salt('bf'))) RETURNING id;`, pgx.NamedArgs{
		"username": user.Username,
		"password": user.Password,
	}).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = errors.Join(err, store.ErrAlreadyExists)
			return
		}
		logger.Log.Warn(err.Error())
	}
	return
}
