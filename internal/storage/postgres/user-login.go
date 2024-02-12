package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (s *storage) UserLogin(ctx context.Context, user models.User) (models.User, error) {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	rows, err := s.pool.Query(ctxQuery, "SELECT id, username, created_at FROM users WHERE username = $1 AND password = crypt($2, password)",
		user.Username, user.Password)
	if err != nil {
		return user, err
	}

	user, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, pgx.ErrTooManyRows) {
			return user, store.ErrNotFound
		}
		return user, err
	}

	return user, nil
}
