package order

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/app/config"
)

func (ts *PostgresTestSuite) TestUpdateStatus() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	userID, err := ts.addTestUser(ctx)
	ts.Require().NoError(err)

	orderData := models.Order{
		UserID: userID,
		Num:    "0866150147",
		Status: models.OrderStatusNew,
	}
	orderID, err := ts.storage.Add(ctx, orderData)
	ts.NoError(err)

	ts.Run("Status Processing", func() {
		err = ts.storage.UpdateStatus(ctx, orderID, models.OrderStatusProcessing)
		ts.NoError(err)

		var orderStatus string
		row := ts.storage.pool.QueryRow(ctx, "SELECT status FROM orders WHERE id = @id;", pgx.NamedArgs{
			"id": orderID,
		})
		err = row.Scan(&orderStatus)
		ts.Equal(models.OrderStatusProcessing.String(), orderStatus)
	})

	ts.Run("Wrong status", func() {
		err = ts.storage.UpdateStatus(ctx, orderID, "test")
		ts.Error(err)
	})

	ts.Run("Undefined status", func() {
		err = ts.storage.UpdateStatus(ctx, orderID, "")
		ts.Error(err)
	})
}

func (ts *PostgresTestSuite) TestAddAccrual() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	userID, err := ts.addTestUser(ctx)
	ts.Require().NoError(err)

	order := models.Order{
		UserID: userID,
		Num:    "0866150147",
		Status: models.OrderStatusNew,
	}
	orderID, err := ts.storage.Add(ctx, order)
	ts.NoError(err)
	order.ID = orderID

	accrual := decimal.NewFromFloatWithExponent(174.682, -config.DecimalExponent)

	type LoyaltyBalanceResult struct {
		UserID    models.UserID   `db:"user_id"`   // ID пользователя
		Current   decimal.Decimal `db:"current"`   // Текущий баланс
		Accrued   decimal.Decimal `db:"accrued"`   // Начислено за всё время
		Withdrawn decimal.Decimal `db:"withdrawn"` // Списано за всё время
	}

	type LoyaltyHistoryResult struct {
		UserID      models.UserID   `db:"user_id"`      // ID пользователя
		OrderNum    string          `db:"order_num"`    // Номер заказа
		Accrual     decimal.Decimal `db:"accrual"`      // Начисление
		Withdrawal  decimal.Decimal `db:"withdrawal"`   // Списание
		ProcessedAt time.Time       `db:"processed_at"` // Дата зачисления или списания
	}

	ts.Run("Positive case", func() {
		err = ts.storage.AddAccrual(ctx, orderID, accrual)
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
