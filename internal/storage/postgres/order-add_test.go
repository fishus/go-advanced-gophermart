package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (ts *PostgresTestSuite) TestOrderAdd() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	userID, err := ts.addTestUser(ctx)
	ts.Require().NoError(err)

	orderData := models.Order{
		UserID: userID,
		Num:    "0866150147",
		Status: models.OrderStatusNew,
	}

	ts.Run("Positive case", func() {
		orderID, err := ts.storage.OrderAdd(ctx, orderData)
		ts.NoError(err)

		var want struct {
			id         string
			userID     string
			num        string
			accrual    float64
			status     string
			uploadedAt time.Time
			updatedAt  time.Time
		}
		err = ts.pool.QueryRow(ctx, "SELECT id, user_id, num, accrual, status, uploaded_at, updated_at FROM orders WHERE id = @id;",
			pgx.NamedArgs{"id": orderID}).Scan(&want.id, &want.userID, &want.num, &want.accrual, &want.status, &want.uploadedAt, &want.updatedAt)
		ts.NoError(err)
		ts.Equal(orderID.String(), want.id)
		ts.Equal(orderData.UserID.String(), want.userID)
		ts.Equal(orderData.Num, want.num)
		ts.Equal(float64(0), want.accrual)
		ts.Equal(orderData.Status.String(), want.status)
		ts.Equal(time.Now().UTC().Round(time.Minute), want.uploadedAt.Round(time.Minute))
		ts.Equal(time.Now().UTC().Round(time.Minute), want.updatedAt.Round(time.Minute))
	})

	ts.Run("DuplicateOrder", func() {
		_, err = ts.storage.OrderAdd(ctx, orderData)
		ts.Error(err)
		ts.ErrorIs(err, store.ErrAlreadyExists)
	})

	ts.Run("IncorrectOrder", func() {
		orderData := models.Order{
			UserID: "",
			Num:    "",
			Status: models.OrderStatusNew,
		}
		_, err := ts.storage.OrderAdd(ctx, orderData)
		ts.Error(err)
		ts.ErrorIs(err, store.ErrIncorrectData)
	})
}
