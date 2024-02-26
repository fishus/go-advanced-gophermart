package postgres

import (
	"context"
	"errors"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
	"github.com/jackc/pgx/v5"
)

func (s *storage) OrderUpdateStatus(ctx context.Context, id models.OrderID, status models.OrderStatus) error {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	_, err := s.pool.Exec(ctxQuery, `UPDATE orders SET status = @status WHERE id = @id;`, pgx.NamedArgs{
		"id":     id.String(),
		"status": status.String(),
	})

	return err
}

func (s *storage) OrderAddAccrual(ctx context.Context, order models.Order, accrual float64) error {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	tx, err := s.pool.Begin(ctxQuery)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctxQuery)

	_, err = tx.Exec(ctxQuery, `UPDATE orders SET accrual = @accrual, status = @status WHERE id = @id;`, pgx.NamedArgs{
		"id":      order.ID.String(),
		"accrual": accrual,
		"status":  models.OrderStatusProcessed,
	})
	if err != nil {
		if errR := tx.Rollback(ctxQuery); errR != nil {
			return errors.Join(err, errR)
		}
		return err
	}

	lh := models.LoyaltyHistory{
		UserID:     order.UserID,
		OrderNum:   order.Num,
		Accrual:    accrual,
		Withdrawal: 0,
	}

	err = s.LoyaltyHistoryAdd(ctx, tx, lh)
	if err != nil {
		if errR := tx.Rollback(ctxQuery); errR != nil {
			return errors.Join(err, errR)
		}
		return err
	}

	lb := models.LoyaltyBalance{
		UserID:    order.UserID,
		Accrued:   accrual,
		Withdrawn: 0,
	}

	err = s.LoyaltyBalanceUpdate(ctx, tx, lb)
	if err != nil {
		if errR := tx.Rollback(ctxQuery); errR != nil {
			return errors.Join(err, errR)
		}
		return err
	}

	return tx.Commit(ctxQuery)
}
