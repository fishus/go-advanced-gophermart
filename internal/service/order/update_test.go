package order

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (ts *OrderServiceTestSuite) TestUpdateStatus() {
	ctx := context.Background()

	orderID := models.OrderID(uuid.New().String())

	testCases := []struct {
		name    string
		status  models.OrderStatus
		wantErr bool
	}{
		{
			"Status Processing",
			models.OrderStatusProcessing,
			false,
		},
		{
			"Status Invalid",
			models.OrderStatusInvalid,
			false,
		},
		{
			"Status Processed",
			models.OrderStatusProcessed,
			false,
		},
		{
			"Status New",
			models.OrderStatusNew,
			false,
		},
		{
			"Wrong status",
			"test1234",
			true,
		},
		{
			"Undefined status",
			models.OrderStatusUndefined,
			true,
		},
	}

	for _, tc := range testCases {
		ts.Run(tc.name, func() {
			mockCall := ts.storage.On("OrderUpdateStatus", ctx, orderID, tc.status).Return(nil)
			err := ts.service.UpdateStatus(ctx, orderID, tc.status)
			if tc.wantErr {
				ts.Error(err)
			} else {
				ts.storage.AssertExpectations(ts.T())
				ts.NoError(err)
			}
			mockCall.Unset()
		})
	}
}

func (ts *OrderServiceTestSuite) TestAddAccrual() {
	ctx := context.Background()

	userID := models.UserID(uuid.New().String())
	orderID := models.OrderID(uuid.New().String())

	mockOrder := models.Order{
		ID:         orderID,
		UserID:     userID,
		Num:        "9890896385",
		Accrual:    0,
		Status:     models.OrderStatusNew,
		UploadedAt: time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	ts.Run("Positive case", func() {
		accrual := 123.456
		mockCallOrderByID := ts.storage.On("OrderByID", ctx, orderID).Return(mockOrder, nil)
		mockCall := ts.storage.On("OrderAddAccrual", ctx, orderID, accrual).Return(nil)
		err := ts.service.AddAccrual(ctx, orderID, accrual)
		ts.NoError(err)
		ts.storage.AssertExpectations(ts.T())
		mockCall.Unset()
		mockCallOrderByID.Unset()
	})

	ts.Run("Negative accrual", func() {
		accrual := -100.0
		err := ts.service.AddAccrual(ctx, orderID, accrual)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrIncorrectData)
		ts.storage.AssertNotCalled(ts.T(), "OrderByID")
		ts.storage.AssertNotCalled(ts.T(), "OrderAddAccrual")
	})

	ts.Run("Status Processed", func() {
		mockOrder.Status = models.OrderStatusProcessed
		accrual := 123.456
		mockCallOrderByID := ts.storage.On("OrderByID", ctx, orderID).Return(mockOrder, nil)
		err := ts.service.AddAccrual(ctx, orderID, accrual)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrOrderRewardReceived)
		ts.storage.AssertExpectations(ts.T())
		ts.storage.AssertNotCalled(ts.T(), "OrderAddAccrual")
		mockCallOrderByID.Unset()
	})

	ts.Run("Order not found", func() {
		accrual := 123.456
		mockCallOrderByID := ts.storage.On("OrderByID", ctx, orderID).Return(models.Order{}, store.ErrNotFound)
		err := ts.service.AddAccrual(ctx, orderID, accrual)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrOrderNotFound)
		ts.storage.AssertExpectations(ts.T())
		ts.storage.AssertNotCalled(ts.T(), "OrderAddAccrual")
		mockCallOrderByID.Unset()
	})
}
