package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/app/config"
	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	stMocks "github.com/fishus/go-advanced-gophermart/internal/storage/mocks"
)

func (ts *OrderServiceTestSuite) TestAddAccrual() {
	ctx := context.Background()

	userID := models.UserID(uuid.New().String())
	orderID := models.OrderID(uuid.New().String())

	mockOrder := models.Order{
		ID:         orderID,
		UserID:     userID,
		Num:        "9890896385",
		Accrual:    decimal.NewFromFloat(0),
		Status:     models.OrderStatusNew,
		UploadedAt: time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	ts.Run("Positive case", func() {
		accrual := decimal.NewFromFloatWithExponent(123.456, -config.DecimalExponent)

		stOrder := stMocks.NewOrderer(ts.T())
		stOrder.EXPECT().GetByID(ctx, orderID).Return(mockOrder, nil)
		stOrder.EXPECT().AddAccrual(ctx, orderID, accrual).Return(nil)
		ts.setStorage(stOrder, nil, nil)

		err := ts.service.AddAccrual(ctx, orderID, accrual)
		ts.NoError(err)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("Negative accrual", func() {
		accrual := decimal.NewFromFloatWithExponent(-100.0, -config.DecimalExponent)
		err := ts.service.AddAccrual(ctx, orderID, accrual)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrIncorrectData)
	})

	ts.Run("Status Processed", func() {
		mockOrder.Status = models.OrderStatusProcessed
		accrual := decimal.NewFromFloatWithExponent(123.456, -config.DecimalExponent)

		stOrder := stMocks.NewOrderer(ts.T())
		stOrder.EXPECT().GetByID(ctx, orderID).Return(mockOrder, nil)
		ts.setStorage(stOrder, nil, nil)

		err := ts.service.AddAccrual(ctx, orderID, accrual)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrOrderRewardReceived)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("Order not found", func() {
		accrual := decimal.NewFromFloatWithExponent(123.456, -config.DecimalExponent)

		stOrder := stMocks.NewOrderer(ts.T())
		stOrder.EXPECT().GetByID(ctx, orderID).Return(models.Order{}, store.ErrNotFound)
		ts.setStorage(stOrder, nil, nil)

		err := ts.service.AddAccrual(ctx, orderID, accrual)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrOrderNotFound)
		ts.storage.AssertExpectations(ts.T())
	})
}
