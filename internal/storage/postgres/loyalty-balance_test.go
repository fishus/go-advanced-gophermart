package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (ts *PostgresTestSuite) TestLoyaltyBalanceUpdate() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	// Setup test data
	userID, err := ts.addTestUser(ctx)
	ts.Require().NoError(err)

	ts.Run("Accrual", func() {
		tx, err := ts.storage.pool.Begin(ctx)
		ts.Require().NoError(err)

		bal := models.LoyaltyBalance{
			UserID:    userID,
			Accrued:   decimal.NewFromFloatWithExponent(768.978, -5),
			Withdrawn: decimal.NewFromFloat(0),
		}
		err = ts.storage.loyaltyBalanceUpdate(ctx, tx, bal)
		ts.NoError(err)
		tx.Commit(ctx)

		want := LoyaltyBalanceResult(bal)
		want.Current = want.Accrued.Sub(want.Withdrawn)

		rows, err := ts.storage.pool.Query(ctx, "SELECT * FROM loyalty_balance WHERE user_id = @userID;", pgx.NamedArgs{
			"userID": userID.String(),
		})
		ts.NoError(err)
		balance, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[LoyaltyBalanceResult])
		ts.NoError(err)
		ts.EqualValues(want, balance)
	})

	ts.Run("Withdraw", func() {
		tx, err := ts.storage.pool.Begin(ctx)
		ts.Require().NoError(err)

		bal := models.LoyaltyBalance{
			UserID:    userID,
			Accrued:   decimal.NewFromFloat(0),
			Withdrawn: decimal.NewFromFloatWithExponent(321.473, -5),
		}
		err = ts.storage.loyaltyBalanceUpdate(ctx, tx, bal)
		ts.NoError(err)
		tx.Commit(ctx)

		want := LoyaltyBalanceResult(bal)
		want.Accrued = decimal.NewFromFloatWithExponent(768.978, -5)
		want.Current = want.Accrued.Sub(want.Withdrawn)

		rows, err := ts.storage.pool.Query(ctx, "SELECT * FROM loyalty_balance WHERE user_id = @userID;", pgx.NamedArgs{
			"userID": userID.String(),
		})
		ts.NoError(err)
		balance, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[LoyaltyBalanceResult])
		ts.NoError(err)
		ts.EqualValues(want, balance)
	})
}

func (ts *PostgresTestSuite) TestLoyaltyAddWithdraw() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	// Setup test data
	userID, err := ts.addTestUser(ctx)
	ts.Require().NoError(err)

	wantBalance := LoyaltyBalanceResult{
		UserID:    userID,
		Accrued:   decimal.NewFromFloatWithExponent(115.387, -5),
		Withdrawn: decimal.NewFromFloat(0),
	}
	wantBalance.Current = wantBalance.Accrued.Sub(wantBalance.Withdrawn)
	_, err = ts.storage.pool.Exec(ctx, `INSERT INTO loyalty_balance (user_id, current, accrued, withdrawn) VALUES (@userID, @current, @accrued, @withdrawn);`, pgx.NamedArgs{
		"userID":    userID,
		"current":   wantBalance.Current,
		"accrued":   wantBalance.Accrued,
		"withdrawn": wantBalance.Withdrawn,
	})
	ts.Require().NoError(err)

	wantBalance.Withdrawn = decimal.NewFromFloatWithExponent(99.995, -5)
	wantBalance.Current = wantBalance.Accrued.Sub(wantBalance.Withdrawn)

	orderNum := "7334280935"

	ts.Run("Positive case", func() {
		wantHistory := LoyaltyHistoryResult{
			UserID:      userID,
			OrderNum:    orderNum,
			Accrual:     decimal.NewFromFloat(0),
			Withdrawal:  wantBalance.Withdrawn,
			ProcessedAt: time.Now().UTC().Round(time.Minute),
		}
		wantHistoryList := make([]LoyaltyHistoryResult, 0, 1)
		wantHistoryList = append(wantHistoryList, wantHistory)

		err = ts.storage.LoyaltyAddWithdraw(ctx, userID, wantHistory.OrderNum, wantBalance.Withdrawn)
		ts.NoError(err)

		rows, err := ts.storage.pool.Query(ctx, "SELECT * FROM loyalty_balance WHERE user_id = @userID;", pgx.NamedArgs{
			"userID": userID.String(),
		})
		ts.NoError(err)

		balance, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[LoyaltyBalanceResult])
		ts.NoError(err)
		ts.EqualValues(wantBalance, balance)

		rows, err = ts.storage.pool.Query(ctx, "SELECT * FROM loyalty_history WHERE user_id = @userID;", pgx.NamedArgs{
			"userID": userID.String(),
		})
		ts.NoError(err)

		history, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[LoyaltyHistoryResult])
		ts.NoError(err)

		for i, h := range history {
			history[i].ProcessedAt = h.ProcessedAt.UTC().Round(time.Minute)
		}

		ts.EqualValues(wantHistoryList, history)
	})

	ts.Run("Low balance", func() {
		err = ts.storage.LoyaltyAddWithdraw(ctx, userID, orderNum, decimal.NewFromFloatWithExponent(959.347, -5))
		ts.Error(err)
		ts.ErrorIs(err, store.ErrLowBalance)
	})
}

func (ts *PostgresTestSuite) TestLoyaltyBalanceByUser() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	// Setup test data
	userID, err := ts.addTestUser(ctx)
	ts.Require().NoError(err)

	want := models.LoyaltyBalance{
		UserID:    userID,
		Accrued:   decimal.NewFromFloatWithExponent(768.978, -5),
		Withdrawn: decimal.NewFromFloatWithExponent(321.473, -5),
	}
	want.Current = want.Accrued.Sub(want.Withdrawn)
	_, err = ts.storage.pool.Exec(ctx, `INSERT INTO loyalty_balance (user_id, current, accrued, withdrawn) VALUES (@userID, @current, @accrued, @withdrawn);`, pgx.NamedArgs{
		"userID":    userID,
		"current":   want.Current,
		"accrued":   want.Accrued,
		"withdrawn": want.Withdrawn,
	})
	ts.Require().NoError(err)

	ts.Run("Positive case", func() {
		balance, err := ts.storage.LoyaltyBalanceByUser(ctx, userID)
		ts.NoError(err)
		ts.EqualValues(want, balance)
	})
}
