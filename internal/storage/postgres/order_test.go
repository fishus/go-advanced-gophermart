package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (ts *PostgresTestSuite) TestOrderByID() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	ts.Run("Return order by ID", func() {
		userID, err := ts.addTestUser(ctx)
		ts.Require().NoError(err)

		orderData := models.Order{
			UserID:     userID,
			Num:        "8020122696",
			Accrual:    0,
			Status:     models.OrderStatusNew,
			UploadedAt: time.Now().UTC().Round(time.Minute),
			UpdatedAt:  time.Now().UTC().Round(time.Minute),
		}

		orderID, err := ts.storage.OrderAdd(ctx, orderData)
		ts.Require().NoError(err)
		orderData.ID = orderID

		order, err := ts.storage.OrderByID(ctx, orderData.ID)
		ts.NoError(err)
		order.UploadedAt = order.UploadedAt.UTC().Round(time.Minute)
		order.UpdatedAt = order.UpdatedAt.UTC().Round(time.Minute)
		ts.EqualValues(orderData, order)
	})

	ts.Run("Order not found", func() {
		orderID := models.OrderID(uuid.New().String())
		_, err := ts.storage.OrderByID(ctx, orderID)
		ts.Error(err)
		ts.ErrorIs(err, store.ErrNotFound)
	})
}

func (ts *PostgresTestSuite) TestOrderByFilter() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	orderNums := []string{"5431720977", "5882492415"}
	orderData := make([]models.Order, len(orderNums))
	for i, orderNum := range orderNums {
		userID, err := ts.addTestUser(ctx)
		ts.Require().NoError(err)

		orderData[i] = models.Order{
			UserID:     userID,
			Num:        orderNum,
			Accrual:    0,
			Status:     models.OrderStatusNew,
			UploadedAt: time.Now().UTC().Round(time.Minute),
			UpdatedAt:  time.Now().UTC().Round(time.Minute),
		}
		orderID, err := ts.storage.OrderAdd(ctx, orderData[i])
		ts.Require().NoError(err)
		orderData[i].ID = orderID
	}

	ts.Run("WithOrderNum", func() {
		order, err := ts.storage.OrderByFilter(ctx, store.WithOrderNum(orderData[0].Num))
		ts.NoError(err)
		order.UploadedAt = order.UploadedAt.UTC().Round(time.Minute)
		order.UpdatedAt = order.UpdatedAt.UTC().Round(time.Minute)
		ts.EqualValues(orderData[0], order)
	})

	ts.Run("WithOrderUserID", func() {
		order, err := ts.storage.OrderByFilter(ctx, store.WithOrderUserID(orderData[1].UserID))
		ts.NoError(err)
		order.UploadedAt = order.UploadedAt.UTC().Round(time.Minute)
		order.UpdatedAt = order.UpdatedAt.UTC().Round(time.Minute)
		ts.EqualValues(orderData[1], order)
	})
}
