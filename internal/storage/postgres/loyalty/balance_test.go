package loyalty

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/app/config"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (ts *PostgresTestSuite) TestBalanceUpdate() {
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
			Accrued:   decimal.NewFromFloatWithExponent(768.978, -config.DecimalExponent),
			Withdrawn: decimal.NewFromFloat(0),
		}
		err = ts.storage.BalanceUpdate(ctx, tx, bal)
		ts.NoError(err)
		tx.Commit(ctx)

		want := BalanceResult(bal)
		want.Current = want.Accrued.Sub(want.Withdrawn)

		rows, err := ts.storage.pool.Query(ctx, "SELECT * FROM loyalty_balance WHERE user_id = @userID;", pgx.NamedArgs{
			"userID": userID.String(),
		})
		ts.NoError(err)
		balance, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[BalanceResult])
		ts.NoError(err)
		ts.EqualValues(want, balance)
	})

	ts.Run("Withdraw", func() {
		tx, err := ts.storage.pool.Begin(ctx)
		ts.Require().NoError(err)

		bal := models.LoyaltyBalance{
			UserID:    userID,
			Accrued:   decimal.NewFromFloat(0),
			Withdrawn: decimal.NewFromFloatWithExponent(321.473, -config.DecimalExponent),
		}
		err = ts.storage.BalanceUpdate(ctx, tx, bal)
		ts.NoError(err)
		tx.Commit(ctx)

		want := BalanceResult(bal)
		want.Accrued = decimal.NewFromFloatWithExponent(768.978, -config.DecimalExponent)
		want.Current = want.Accrued.Sub(want.Withdrawn)

		rows, err := ts.storage.pool.Query(ctx, "SELECT * FROM loyalty_balance WHERE user_id = @userID;", pgx.NamedArgs{
			"userID": userID.String(),
		})
		ts.NoError(err)
		balance, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[BalanceResult])
		ts.NoError(err)
		ts.EqualValues(want, balance)
	})
}

func (ts *PostgresTestSuite) TestAddWithdraw() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	// Setup test data
	userID, err := ts.addTestUser(ctx)
	ts.Require().NoError(err)

	wantBalance := BalanceResult{
		UserID:    userID,
		Accrued:   decimal.NewFromFloatWithExponent(115.387, -config.DecimalExponent),
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

	wantBalance.Withdrawn = decimal.NewFromFloatWithExponent(99.995, -config.DecimalExponent)
	wantBalance.Current = wantBalance.Accrued.Sub(wantBalance.Withdrawn)

	orderNum := "7334280935"

	ts.Run("Positive case", func() {
		wantHistory := HistoryResult{
			UserID:      userID,
			OrderNum:    orderNum,
			Accrual:     decimal.NewFromFloat(0),
			Withdrawal:  wantBalance.Withdrawn,
			ProcessedAt: time.Now().UTC().Round(time.Minute),
		}
		wantHistoryList := make([]HistoryResult, 0, 1)
		wantHistoryList = append(wantHistoryList, wantHistory)

		err = ts.storage.AddWithdraw(ctx, userID, wantHistory.OrderNum, wantBalance.Withdrawn)
		ts.NoError(err)

		rows, err := ts.storage.pool.Query(ctx, "SELECT * FROM loyalty_balance WHERE user_id = @userID;", pgx.NamedArgs{
			"userID": userID.String(),
		})
		ts.NoError(err)

		balance, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[BalanceResult])
		ts.NoError(err)
		ts.EqualValues(wantBalance, balance)

		rows, err = ts.storage.pool.Query(ctx, "SELECT * FROM loyalty_history WHERE user_id = @userID;", pgx.NamedArgs{
			"userID": userID.String(),
		})
		ts.NoError(err)

		history, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[HistoryResult])
		ts.NoError(err)

		for i, h := range history {
			history[i].ProcessedAt = h.ProcessedAt.UTC().Round(time.Minute)
		}

		ts.EqualValues(wantHistoryList, history)
	})

	ts.Run("Low balance", func() {
		err = ts.storage.AddWithdraw(ctx, userID, orderNum, decimal.NewFromFloatWithExponent(959.347, -config.DecimalExponent))
		ts.Error(err)
		ts.ErrorIs(err, store.ErrLowBalance)
	})
}

func (ts *PostgresTestSuite) TestBalanceByUser() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	// Setup test data
	userID, err := ts.addTestUser(ctx)
	ts.Require().NoError(err)

	want := models.LoyaltyBalance{
		UserID:    userID,
		Accrued:   decimal.NewFromFloatWithExponent(768.978, -config.DecimalExponent),
		Withdrawn: decimal.NewFromFloatWithExponent(321.473, -config.DecimalExponent),
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
		balance, err := ts.storage.BalanceByUser(ctx, userID)
		ts.NoError(err)
		ts.EqualValues(want, balance)
	})
}
