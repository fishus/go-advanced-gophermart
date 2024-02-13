package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (s *storage) UserLogin(ctx context.Context, user models.User) (models.UserID, error) {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	var userID models.UserID
	row := s.pool.QueryRow(ctxQuery, "SELECT id FROM users WHERE username = @username AND password = crypt(@password, password);",
		pgx.NamedArgs{"username": user.Username, "password": user.Password})
	err := row.Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return userID, store.ErrNotFound
		}
		logger.Log.Warn(err.Error())
		return userID, err
	}

	return userID, nil
}
