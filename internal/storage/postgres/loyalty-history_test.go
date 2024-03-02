package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (ts *PostgresTestSuite) TestLoyaltyHistoryAdd() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	// Setup test data
	userID, err := ts.addTestUser(ctx)
	ts.Require().NoError(err)

	ts.Run("Positive case", func() {
		tx, err := ts.storage.pool.Begin(ctx)
		ts.Require().NoError(err)

		wantHistory := make([]models.LoyaltyHistory, 2)
		wantHistory[0] = models.LoyaltyHistory{
			UserID:      userID,
			OrderNum:    "5347676263",
			Accrual:     123.456,
			Withdrawal:  0,
			ProcessedAt: time.Now().UTC().Round(time.Minute),
		}
		wantHistory[1] = models.LoyaltyHistory{
			UserID:      userID,
			OrderNum:    "8163091187",
			Accrual:     0,
			Withdrawal:  654.321,
			ProcessedAt: time.Now().UTC().Round(time.Minute),
		}

		for _, h := range wantHistory {
			err = ts.storage.loyaltyHistoryAdd(ctx, tx, h)
			ts.NoError(err)
		}
		tx.Commit(ctx)

		rows, err := ts.storage.pool.Query(ctx, "SELECT * FROM loyalty_history WHERE user_id = @userID;", pgx.NamedArgs{
			"userID": userID.String(),
		})
		ts.NoError(err)
		historyData, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[LoyaltyHistoryResult])
		ts.NoError(err)
		historyList := make([]models.LoyaltyHistory, 0)
		for _, h := range historyData {
			history := models.LoyaltyHistory(h)
			history.ProcessedAt = h.ProcessedAt.UTC().Round(time.Minute)
			historyList = append(historyList, history)
		}
		ts.EqualValues(wantHistory, historyList)
	})
}

func (ts *PostgresTestSuite) TestLoyaltyHistoryByUser() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	// Setup test data
	userID, err := ts.addTestUser(ctx)
	ts.Require().NoError(err)

	wantHistory := make([]models.LoyaltyHistory, 2)
	wantHistory[0] = models.LoyaltyHistory{
		UserID:      userID,
		OrderNum:    "6825296715",
		Accrual:     123.456,
		Withdrawal:  0,
		ProcessedAt: time.Now().UTC().Round(time.Minute),
	}
	wantHistory[1] = models.LoyaltyHistory{
		UserID:      userID,
		OrderNum:    "8215993786",
		Accrual:     0,
		Withdrawal:  654.321,
		ProcessedAt: time.Now().UTC().Round(time.Minute),
	}

	for _, h := range wantHistory {
		_, err = ts.storage.pool.Exec(ctx, `INSERT INTO loyalty_history (user_id, order_num, accrual, withdrawal) VALUES (@userID, @orderNum, @accrual, @withdrawal);`, pgx.NamedArgs{
			"userID":     h.UserID,
			"orderNum":   h.OrderNum,
			"accrual":    h.Accrual,
			"withdrawal": h.Withdrawal,
		})
		ts.Require().NoError(err)
	}

	ts.Run("Positive case", func() {
		history, err := ts.storage.LoyaltyHistoryByUser(ctx, userID)
		for i, h := range history {
			history[i].ProcessedAt = h.ProcessedAt.UTC().Round(time.Minute)
		}
		ts.NoError(err)
		ts.EqualValues(wantHistory, history)
	})
}