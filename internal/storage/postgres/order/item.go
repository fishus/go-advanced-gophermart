package order

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (s *storage) GetByID(ctx context.Context, id models.OrderID) (order models.Order, err error) {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	rows, err := s.pool.Query(ctxQuery, "SELECT * FROM orders WHERE id = @id;", pgx.NamedArgs{
		"id": id,
	})
	if err != nil {
		return
	}

	orderResult, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[OrderResult])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, pgx.ErrTooManyRows) {
			err = store.ErrNotFound
			return
		}
		logger.Log.Warn(err.Error())
		return
	}

	order = models.Order(orderResult)
	return
}

func (s *storage) txGetByIDForUpdate(ctx context.Context, tx pgx.Tx, id models.OrderID) (order models.Order, err error) {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	rows, err := tx.Query(ctxQuery, "SELECT * FROM orders WHERE id = @id FOR UPDATE;", pgx.NamedArgs{
		"id": id,
	})
	if err != nil {
		return
	}

	orderResult, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[OrderResult])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, pgx.ErrTooManyRows) {
			err = store.ErrNotFound
			return
		}
		logger.Log.Warn(err.Error())
		return
	}

	order = models.Order(orderResult)
	return
}

func (s *storage) GetByFilter(ctx context.Context, filters ...store.OrderFilter) (order models.Order, err error) {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	f := &store.OrderFilters{}
	for _, filter := range filters {
		filter(f)
	}
	if f.IsEmpty() {
		err = errors.New("at least one filter required")
		return
	}

	queryFilter := make([]string, 0)
	namedArgs := pgx.NamedArgs{}

	if f.Num != "" {
		queryFilter = append(queryFilter, `num = @num`)
		namedArgs["num"] = f.Num
	}

	if f.UserID != "" {
		queryFilter = append(queryFilter, `user_id = @userID`)
		namedArgs["userID"] = f.UserID
	}

	filterStr := strings.Join(queryFilter, ` AND `)

	rows, err := s.pool.Query(ctxQuery, `SELECT * FROM orders WHERE `+filterStr+` LIMIT 1;`, namedArgs)
	if err != nil {
		return
	}

	orderResult, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[OrderResult])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = store.ErrNotFound
			return
		}
		logger.Log.Warn(err.Error())
		return
	}

	order = models.Order(orderResult)
	return
}
