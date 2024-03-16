package order

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

func (s *storage) Add(ctx context.Context, order models.Order) (orderID models.OrderID, err error) {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	if order.UserID == "" || order.Num == "" {
		err = store.ErrIncorrectData
		return
	}

	err = s.pool.QueryRow(ctxQuery, `INSERT INTO orders (user_id, num, status) VALUES (@userID, @num, @status) RETURNING id;`,
		pgx.NamedArgs{
			"userID": order.UserID,
			"num":    order.Num,
			"status": order.Status,
		}).Scan(&orderID)
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
