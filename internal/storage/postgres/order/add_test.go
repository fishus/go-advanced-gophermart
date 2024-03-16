package order

import (
	"context"
	"time"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (ts *PostgresTestSuite) TestAdd() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	userID, err := ts.addTestUser(ctx)
	ts.Require().NoError(err)

	data := models.Order{
		UserID:     userID,
		Num:        "0866150147",
		Accrual:    decimal.NewFromFloat(0),
		Status:     models.OrderStatusNew,
		UploadedAt: time.Now().UTC().Round(time.Minute),
		UpdatedAt:  time.Now().UTC().Round(time.Minute),
	}

	ts.Run("Positive case", func() {
		orderID, err := ts.storage.Add(ctx, data)
		ts.NoError(err)
		data.ID = orderID

		var want OrderResult
		err = ts.pool.QueryRow(ctx, "SELECT id, user_id, num, accrual, status, uploaded_at, updated_at FROM orders WHERE id = @id;",
			pgx.NamedArgs{"id": orderID}).Scan(&want.ID, &want.UserID, &want.Num, &want.Accrual, &want.Status, &want.UploadedAt, &want.UpdatedAt)
		ts.NoError(err)
		want.UploadedAt = want.UploadedAt.UTC().Round(time.Minute)
		want.UpdatedAt = want.UpdatedAt.UTC().Round(time.Minute)
		ts.EqualValues(OrderResult(data), want)
	})

	ts.Run("DuplicateOrder", func() {
		_, err = ts.storage.Add(ctx, data)
		ts.Error(err)
		ts.ErrorIs(err, store.ErrAlreadyExists)
	})

	ts.Run("IncorrectOrder", func() {
		data := models.Order{
			UserID: "",
			Num:    "",
			Status: models.OrderStatusNew,
		}
		_, err := ts.storage.Add(ctx, data)
		ts.Error(err)
		ts.ErrorIs(err, store.ErrIncorrectData)
	})
}
