package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (s *storage) loyaltyBalanceUpdate(ctx context.Context, tx pgx.Tx, balance models.LoyaltyBalance) error {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	_, err := tx.Exec(ctxQuery, `
WITH b AS (
	SELECT @accrued::numeric as accrued, @withdrawn::numeric as withdrawn
)
INSERT INTO loyalty_balance SELECT @userID::uuid AS user_id, (accrued - withdrawn) AS current, accrued, withdrawn FROM b
ON CONFLICT (user_id) DO UPDATE SET
accrued = loyalty_balance.accrued + EXCLUDED.accrued,
withdrawn = loyalty_balance.withdrawn + EXCLUDED.withdrawn,
current = ((loyalty_balance.accrued + EXCLUDED.accrued) - (loyalty_balance.withdrawn + EXCLUDED.withdrawn));
`,
		pgx.NamedArgs{
			"userID":    balance.UserID.String(),
			"accrued":   balance.Accrued,
			"withdrawn": balance.Withdrawn,
		})
	return err
}

func (s *storage) LoyaltyAddWithdraw(ctx context.Context, userID models.UserID, orderNum string, withdraw float64) error {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	tx, err := s.pool.Begin(ctxQuery)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctxQuery)

	// Проверка баланса
	balance, err := s.LoyaltyBalanceByUser(ctx, userID)
	if err != nil {
		return err
	}

	if balance.Current < withdraw {
		return store.ErrLowBalance
	}

	lh := models.LoyaltyHistory{
		UserID:     userID,
		OrderNum:   orderNum,
		Accrual:    0,
		Withdrawal: withdraw,
	}

	err = s.loyaltyHistoryAdd(ctx, tx, lh)
	if err != nil {
		if errR := tx.Rollback(ctxQuery); errR != nil {
			return errors.Join(err, errR)
		}
		return err
	}

	lb := models.LoyaltyBalance{
		UserID:    userID,
		Accrued:   0,
		Withdrawn: withdraw,
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

func (s *storage) LoyaltyBalanceByUser(ctx context.Context, userID models.UserID) (balance models.LoyaltyBalance, err error) {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	rows, err := s.pool.Query(ctxQuery, "SELECT * FROM loyalty_balance WHERE user_id = @userID;", pgx.NamedArgs{
		"userID": userID.String(),
	})
	if err != nil {
		return
	}

	balanceResult, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[LoyaltyBalanceResult])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, pgx.ErrTooManyRows) {
			err = store.ErrNotFound
			return
		}
		logger.Log.Warn(err.Error())
		return
	}

	balance = models.LoyaltyBalance(balanceResult)
	return
}
