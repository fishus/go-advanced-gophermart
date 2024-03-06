package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
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

func (s *storage) OrderAddAccrual(ctx context.Context, orderID models.OrderID, accrual decimal.Decimal) error {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	tx, err := s.pool.Begin(ctxQuery)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctxQuery)

	order, err := s.txOrderByID(ctx, tx, orderID)
	if err != nil {
		return err
	}

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
		Withdrawal: decimal.NewFromFloat(0),
	}

	err = s.loyaltyHistoryAdd(ctx, tx, lh)
	if err != nil {
		if errR := tx.Rollback(ctxQuery); errR != nil {
			return errors.Join(err, errR)
		}
		return err
	}

	lb := models.LoyaltyBalance{
		UserID:    order.UserID,
		Accrued:   accrual,
		Withdrawn: decimal.NewFromFloat(0),
	}

	err = s.loyaltyBalanceUpdate(ctx, tx, lb)
	if err != nil {
		if errR := tx.Rollback(ctxQuery); errR != nil {
			return errors.Join(err, errR)
		}
		return err
	}

	return tx.Commit(ctxQuery)
}
