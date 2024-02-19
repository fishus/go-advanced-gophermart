package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
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
