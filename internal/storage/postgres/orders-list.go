package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (s *storage) OrdersByFilter(ctx context.Context, limit int, filters ...store.OrderFilter) ([]models.Order, error) {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	f := &store.OrderFilters{}
	for _, filter := range filters {
		filter(f)
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
	if filterStr != "" {
		filterStr = "WHERE " + filterStr
	}

	limitStr := ""
	if limit > 0 {
		limitStr = fmt.Sprintf("LIMIT %d", limit)
	}

	rows, err := s.pool.Query(ctxQuery, `SELECT * FROM orders `+filterStr+` `+limitStr+`;`, namedArgs)
	if err != nil {
		return nil, err
	}

	orderResult, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[OrderResult])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		logger.Log.Warn(err.Error())
		return nil, err
	}

	return listResultsToOrders(orderResult), nil
}
