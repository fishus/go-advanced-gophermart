package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (s *storage) UserAdd(ctx context.Context, user models.User) (models.UserID, error) {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	var userID models.UserID
	err := s.pool.QueryRow(ctxQuery, `INSERT INTO users (username, password) VALUES (@username, crypt(@password, gen_salt('bf'))) RETURNING id;`,
		pgx.NamedArgs{"username": user.Username, "password": user.Password}).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return userID, errors.Join(err, store.ErrAlreadyExists)
		}
		logger.Log.Warn(err.Error())
		return userID, err
	}

	return userID, nil
}
