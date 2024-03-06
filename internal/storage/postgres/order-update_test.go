package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (ts *PostgresTestSuite) TestOrderUpdateStatus() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	userID, err := ts.addTestUser(ctx)
	ts.Require().NoError(err)

	orderData := models.Order{
		UserID: userID,
		Num:    "0866150147",
		Status: models.OrderStatusNew,
	}
	orderID, err := ts.storage.OrderAdd(ctx, orderData)
	ts.NoError(err)

	ts.Run("Status Processing", func() {
		err = ts.storage.OrderUpdateStatus(ctx, orderID, models.OrderStatusProcessing)
		ts.NoError(err)

		var orderStatus string
		row := ts.storage.pool.QueryRow(ctx, "SELECT status FROM orders WHERE id = @id;", pgx.NamedArgs{
			"id": orderID,
		})
		err = row.Scan(&orderStatus)
		ts.Equal(models.OrderStatusProcessing.String(), orderStatus)
	})

	ts.Run("Wrong status", func() {
		err = ts.storage.OrderUpdateStatus(ctx, orderID, "test")
		ts.Error(err)
	})

	ts.Run("Undefined status", func() {
		err = ts.storage.OrderUpdateStatus(ctx, orderID, "")
		ts.Error(err)
	})
}

func (ts *PostgresTestSuite) TestOrderAddAccrual() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	userID, err := ts.addTestUser(ctx)
	ts.Require().NoError(err)

	order := models.Order{
		UserID: userID,
		Num:    "0866150147",
		Status: models.OrderStatusNew,
	}
	orderID, err := ts.storage.OrderAdd(ctx, order)
	ts.NoError(err)
	order.ID = orderID

	accrual := decimal.NewFromFloatWithExponent(174.682, -5)

	ts.Run("Positive case", func() {
		err = ts.storage.OrderAddAccrual(ctx, orderID, accrual)
		ts.NoError(err)

		// Check updated order
		var (
			orderAccrual decimal.Decimal
			orderStatus  string
		)
		row := ts.storage.pool.QueryRow(ctx, "SELECT accrual, status FROM orders WHERE id = @id;", pgx.NamedArgs{
			"id": orderID,
		})
		err = row.Scan(&orderAccrual, &orderStatus)

		if orderAccrual.IsZero() {
			orderAccrual = decimal.NewFromFloat(0)
		}

		ts.Equal(accrual, orderAccrual)
		ts.Equal(models.OrderStatusProcessed.String(), orderStatus)

		// Check updated balance
		wantBalance := LoyaltyBalanceResult{
			UserID:    userID,
			Accrued:   accrual,
			Withdrawn: decimal.NewFromFloat(0),
		}
		wantBalance.Current = wantBalance.Accrued.Sub(wantBalance.Withdrawn)
		rows, err := ts.storage.pool.Query(ctx, "SELECT * FROM loyalty_balance WHERE user_id = @userID;", pgx.NamedArgs{
			"userID": userID.String(),
		})
		ts.NoError(err)
		balance, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[LoyaltyBalanceResult])
		ts.NoError(err)
		ts.EqualValues(wantBalance, balance)

		// Check updated history
		wantHistory := make([]LoyaltyHistoryResult, 1)
		wantHistory[0] = LoyaltyHistoryResult{
			UserID:      userID,
			OrderNum:    order.Num,
			Accrual:     accrual,
			Withdrawal:  decimal.NewFromFloat(0),
			ProcessedAt: time.Now().UTC().Round(time.Minute),
		}
		rows, err = ts.storage.pool.Query(ctx, "SELECT * FROM loyalty_history WHERE user_id = @userID;", pgx.NamedArgs{
			"userID": userID.String(),
		})
		ts.NoError(err)
		history, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[LoyaltyHistoryResult])
		ts.NoError(err)
		for i, h := range history {
			history[i].ProcessedAt = h.ProcessedAt.UTC().Round(time.Minute)
		}
		ts.EqualValues(wantHistory, history)
	})
}
