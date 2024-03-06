package order

import (
	"context"
	"slices"
	"time"

	"github.com/shopspring/decimal"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (ts *PostgresTestSuite) TestListByFilter() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	userID, err := ts.addTestUser(ctx)
	ts.Require().NoError(err)

	orderData := make([]models.Order, 3)
	orderData[0] = models.Order{
		Num:    "4581347475",
		Status: models.OrderStatusNew,
	}
	orderData[1] = models.Order{
		Num:    "4941161632",
		Status: models.OrderStatusProcessing,
	}
	orderData[2] = models.Order{
		Num:    "6860728705",
		Status: models.OrderStatusProcessing,
	}
	for i := 0; i < len(orderData); i++ {
		orderData[i].UserID = userID
		orderData[i].Accrual = decimal.NewFromFloat(0)
		orderData[i].UploadedAt = time.Now().UTC().Round(time.Minute)
		orderData[i].UpdatedAt = time.Now().UTC().Round(time.Minute)
		orderID, err := ts.storage.Add(ctx, orderData[i])
		ts.Require().NoError(err)
		orderData[i].ID = orderID
	}

	ts.Run("WithOrderNum", func() {
		orders, err := ts.storage.ListByFilter(ctx, 10, store.WithOrderNum(orderData[0].Num))
		ts.NoError(err)
		ts.Equal(orderData[0].Num, orders[0].Num)
		orders[0].UploadedAt = orders[0].UploadedAt.UTC().Round(time.Minute)
		orders[0].UpdatedAt = orders[0].UpdatedAt.UTC().Round(time.Minute)
		ts.EqualValues(orderData[0], orders[0])
	})

	ts.Run("WithOrderUserID", func() {
		limit := 2
		orders, err := ts.storage.ListByFilter(ctx, limit, store.WithOrderUserID(userID))
		ts.NoError(err)
		ts.Equal(limit, len(orders))
		for _, order := range orders {
			i := slices.IndexFunc(orderData, func(o models.Order) bool {
				return o.Num == order.Num
			})
			ts.Equal(userID, order.UserID)
			order.UploadedAt = order.UploadedAt.UTC().Round(time.Minute)
			order.UpdatedAt = order.UpdatedAt.UTC().Round(time.Minute)
			ts.EqualValues(orderData[i], order)
		}
	})

	ts.Run("WithOrderStatus", func() {
		orders, err := ts.storage.ListByFilter(ctx, 10, store.WithOrderStatus(models.OrderStatusNew))
		ts.NoError(err)
		for _, order := range orders {
			i := slices.IndexFunc(orderData, func(o models.Order) bool {
				return o.Num == order.Num
			})
			ts.Equal(models.OrderStatusNew, order.Status)
			order.UploadedAt = order.UploadedAt.UTC().Round(time.Minute)
			order.UpdatedAt = order.UpdatedAt.UTC().Round(time.Minute)
			ts.EqualValues(orderData[i], order)
		}
	})

	ts.Run("WithOrderStatuses", func() {
		orders, err := ts.storage.ListByFilter(ctx, 10, store.WithOrderStatuses(models.OrderStatusProcessing))
		ts.NoError(err)
		for _, order := range orders {
			i := slices.IndexFunc(orderData, func(o models.Order) bool {
				return o.Num == order.Num
			})
			ts.Equal(models.OrderStatusProcessing, order.Status)
			order.UploadedAt = order.UploadedAt.UTC().Round(time.Minute)
			order.UpdatedAt = order.UpdatedAt.UTC().Round(time.Minute)
			ts.EqualValues(orderData[i], order)
		}
	})
}
