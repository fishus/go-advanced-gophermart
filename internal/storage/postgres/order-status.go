package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (s *storage) OrderSetStatusBatch(ctx context.Context, idList []models.OrderID, newStatus models.OrderStatus) error {
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
