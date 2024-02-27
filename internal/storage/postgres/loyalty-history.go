package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (s *storage) loyaltyHistoryAdd(ctx context.Context, tx pgx.Tx, history models.LoyaltyHistory) error {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	_, err := tx.Exec(ctxQuery, `INSERT INTO loyalty_history (user_id, order_num, accrual, withdrawal) VALUES (@userID, @orderNum, @accrual, @withdrawal)`,
		pgx.NamedArgs{
			"userID":     history.UserID.String(),
			"orderNum":   history.OrderNum,
			"accrual":    history.Accrual,
			"withdrawal": history.Withdrawal,
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
	var curBal float64
	row := tx.QueryRow(ctxQuery, "SELECT current FROM loyalty_balance WHERE user_id = @userID;", pgx.NamedArgs{
		"userID": userID,
	})
	err = row.Scan(&curBal)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return store.ErrNotFound
		}
		return err
	}

	if curBal < withdraw {
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

func (s *storage) LoyaltyHistoryByUser(ctx context.Context, userID models.UserID) ([]models.LoyaltyHistory, error) {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	rows, err := s.pool.Query(ctxQuery, `SELECT * FROM loyalty_history WHERE user_id = @userID ORDER BY processed_at ASC;`, pgx.NamedArgs{
		"userID": userID.String(),
	})
	if err != nil {
		return nil, err
	}

	historyResult, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[LoyaltyHistoryResult])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		logger.Log.Warn(err.Error())
		return nil, err
	}

	return listResultsToLoyaltyHistory(historyResult), nil
}
