package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (s *storage) UserAdd(ctx context.Context, user models.User) (models.User, error) {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	var row struct {
		id   string
		time time.Time
	}
	err := s.pool.QueryRow(ctxQuery, `INSERT INTO users (username, password) VALUES (@username, crypt(@password, gen_salt('bf'))) RETURNING id, created_at;`,
		pgx.NamedArgs{"username": user.Username, "password": user.Password}).Scan(&row.id, &row.time)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return user, errors.Join(err, store.ErrAlreadyExists)
		}
		return user, err
	}

	user.ID = models.UserID(row.id)
	user.CreatedAt = row.time

	return user, nil
}
