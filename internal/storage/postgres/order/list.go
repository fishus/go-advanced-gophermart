package order

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

func (s *storage) ListByFilter(ctx context.Context, limit int, filters ...store.OrderFilter) ([]models.Order, error) {
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
		namedArgs["userID"] = f.UserID.String()
	}

	if len(f.Statuses) == 1 {
		queryFilter = append(queryFilter, `status = @status`)
		namedArgs["status"] = f.Statuses[0].String()
	} else if len(f.Statuses) > 1 {
		queryFilter = append(queryFilter, `status = ANY(@status)`)
		statuses := make([]string, len(f.Statuses))
		for i, status := range f.Statuses {
			statuses[i] = status.String()
		}
		namedArgs["status"] = statuses
	}

	if len(f.ID) == 1 {
		queryFilter = append(queryFilter, `id = @id`)
		namedArgs["id"] = f.ID[0].String()
	} else if len(f.ID) > 1 {
		queryFilter = append(queryFilter, `id = ANY(@id)`)
		idList := make([]string, len(f.ID))
		for i, id := range f.ID {
			idList[i] = id.String()
		}
		namedArgs["id"] = idList
	}

	filterStr := strings.Join(queryFilter, ` AND `)
	if filterStr != "" {
		filterStr = "WHERE " + filterStr
	}

	orderBy := make([]string, 0)
	if len(f.OrderBy) == 0 {
		f.OrderBy = append(f.OrderBy, struct {
			Field store.OrderByField
			Dir   store.OrderByDirection
		}{Field: store.OrderByUploadedAt, Dir: store.OrderByAsc})
	}
	for _, o := range f.OrderBy {
		orderBy = append(orderBy, fmt.Sprintf("%s %s", o.Field, o.Dir))
	}
	orderByStr := "ORDER BY " + strings.Join(orderBy, `, `)

	limitStr := ""
	if limit > 0 {
		limitStr = fmt.Sprintf("LIMIT %d", limit)
	}

	rows, err := s.pool.Query(ctxQuery, `SELECT * FROM orders `+filterStr+` `+orderByStr+` `+limitStr+`;`, namedArgs)
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
