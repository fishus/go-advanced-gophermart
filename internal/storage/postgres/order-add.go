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

func (s *storage) OrderAdd(ctx context.Context, order models.Order) (models.OrderID, error) {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	var orderID models.OrderID
	err := s.pool.QueryRow(ctxQuery, `INSERT INTO orders (user_id, num, status) VALUES (@userID, @num, @status) RETURNING id;`,
		pgx.NamedArgs{"userID": order.UserID, "num": order.Num, "status": order.Status}).Scan(&orderID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return orderID, errors.Join(err, store.ErrAlreadyExists)
		}
		logger.Log.Warn(err.Error())
		return orderID, err
	}

	return orderID, nil
}
