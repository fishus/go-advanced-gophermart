package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (s *storage) OrderByID(ctx context.Context, id models.OrderID) (*models.Order, error) {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	rows, err := s.pool.Query(ctxQuery, "SELECT id, user_id, num, accrual, status, uploaded_at FROM orders WHERE id = @id;", pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, err
	}

	orderResult, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[OrderResult])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, pgx.ErrTooManyRows) {
			return nil, store.ErrNotFound
		}
		logger.Log.Warn(err.Error())
		return nil, err
	}

	order := models.Order(orderResult)
	return &order, nil
}

func (s *storage) OrderByFilter(ctx context.Context, filters ...store.OrderFilter) (*models.Order, error) {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	f := &store.OrderFilters{}
	for _, filter := range filters {
		filter(f)
	}
	if f.IsEmpty() {
		return nil, errors.New("at least one filter required")
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

	rows, err := s.pool.Query(ctxQuery, `SELECT id, user_id, num, accrual, status, uploaded_at FROM orders WHERE `+filterStr+` LIMIT 1;`, namedArgs)
	if err != nil {
		return nil, err
	}

	orderResult, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[OrderResult])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		logger.Log.Warn(err.Error())
		return nil, err
	}

	order := models.Order(orderResult)
	return &order, nil
}
