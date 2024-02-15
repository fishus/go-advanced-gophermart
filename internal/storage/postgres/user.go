package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (s *storage) UserByID(ctx context.Context, id models.UserID) (*models.User, error) {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	rows, err := s.pool.Query(ctxQuery, "SELECT id, username, created_at FROM users WHERE id = @id;", pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, err
	}

	userResult, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[UserResult])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, pgx.ErrTooManyRows) {
			return nil, store.ErrNotFound
		}
		logger.Log.Warn(err.Error())
		return nil, err
	}

	user := models.User(userResult)
	return &user, nil
}
