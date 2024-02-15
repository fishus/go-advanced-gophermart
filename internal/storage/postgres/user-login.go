package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (s *storage) UserLogin(ctx context.Context, user models.User) (userID models.UserID, err error) {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	row := s.pool.QueryRow(ctxQuery, "SELECT id FROM users WHERE username = @username AND password = crypt(@password, password);",
		pgx.NamedArgs{"username": user.Username, "password": user.Password})
	err = row.Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = store.ErrNotFound
			return
		}
		logger.Log.Warn(err.Error())
	}
	return
}
