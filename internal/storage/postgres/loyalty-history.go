package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (s *storage) LoyaltyHistoryAdd(ctx context.Context, tx pgx.Tx, history models.LoyaltyHistory) error {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	_, err := tx.Exec(ctxQuery, `INSERT INTO loyalty_history (user_id, order_id, accrual, withdrawal) VALUES (@userID, @orderID, @accrual, @withdrawal)`,
		pgx.NamedArgs{
			"userID":     history.UserID.String(),
			"orderID":    history.OrderID.String(),
			"accrual":    history.Accrual,
			"withdrawal": history.Withdrawal,
		})
	return err
}
