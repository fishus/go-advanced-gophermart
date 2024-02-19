package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

// OrderResetProcessingStatus reset the status of unfinished orders to the status of new orders.
func (s *storage) OrderResetProcessingStatus(ctx context.Context) error {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	_, err := s.pool.Exec(ctxQuery, `UPDATE orders SET status = @statusNew WHERE status = @statusProcessing;`,
		pgx.NamedArgs{
			"statusNew":        models.OrderStatusNew,
			"statusProcessing": models.OrderStatusProcessing,
		})
	if err != nil {
		logger.Log.Warn(err.Error())
	}
	return err
}

// OrderMoveToProcessing selects N new orders and changes their status to "processing" and then returns a list of these orders
func (s *storage) OrderMoveToProcessing(ctx context.Context, limit int) (list []models.Order, err error) {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	tx, err := s.pool.Begin(ctxQuery)
	if err != nil {
		return
	}

	// List of orders before status changes
	list, err = s.OrdersByFilter(ctx, limit, store.WithOrderStatus(models.OrderStatusNew), store.WithOrderBy(store.OrderByUploadedAt, store.OrderByAsc))
	if err != nil {
		tx.Commit(ctxQuery)
		return
	}

	if len(list) == 0 {
		tx.Commit(ctxQuery)
		return
	}

	stmt, err := tx.Prepare(ctxQuery, "",
		`UPDATE orders SET status = $1 WHERE id = $2;`)
	if err != nil {
		if errR := tx.Rollback(ctxQuery); errR != nil {
			err = errors.Join(err, errR)
		}
		return
	}

	var oIDList []models.OrderID
	for _, order := range list {
		oIDList = append(oIDList, order.ID)
		_, err = tx.Exec(ctxQuery, stmt.SQL, models.OrderStatusProcessing, order.ID.String())
		if err != nil {
			if errR := tx.Rollback(ctxQuery); errR != nil {
				err = errors.Join(err, errR)
			}
			return
		}
	}

	tx.Commit(ctxQuery)

	// List of orders after status changes
	list, err = s.OrdersByFilter(ctx, limit, store.WithOrderIDList(oIDList...), store.WithOrderBy(store.OrderByUploadedAt, store.OrderByAsc))
	if err != nil {
		if errR := tx.Rollback(ctxQuery); errR != nil {
			err = errors.Join(err, errR)
		}
		return
	}

	return
}

func (s *storage) OrderSetStatus(ctx context.Context, idList []models.OrderID, newStatus models.OrderStatus) error {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	if len(idList) == 0 {
		return nil
	}

	tx, err := s.pool.Begin(ctxQuery)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctxQuery)

	batch := &pgx.Batch{}
	for _, id := range idList {
		args := pgx.NamedArgs{
			"id":     id.String(),
			"status": newStatus.String(),
		}
		batch.Queue(`UPDATE orders SET status = @status WHERE id = @id;`, args)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	for range idList {
		_, err = results.Exec()
		if err != nil {
			return err
		}
	}
	err = results.Close()
	if err != nil {
		return err
	}

	return tx.Commit(ctxQuery)
}
