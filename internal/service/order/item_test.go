package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	stMocks "github.com/fishus/go-advanced-gophermart/internal/storage/mocks"
)

func (ts *OrderServiceTestSuite) TestGetByID() {
	ctx := context.Background()

	orderID := models.OrderID(uuid.New().String())

	ts.Run("Positive case", func() {
		want := models.Order{
			ID:         orderID,
			UserID:     models.UserID(uuid.New().String()),
			Num:        "9400781309",
			Accrual:    decimal.NewFromFloat(0),
			Status:     models.OrderStatusNew,
			UploadedAt: time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		}

		stOrder := stMocks.NewOrderer(ts.T())
		stOrder.EXPECT().GetByID(ctx, orderID).Return(want, nil)
		ts.setStorage(stOrder, nil, nil)

		list, err := ts.service.GetByID(ctx, orderID)
		ts.NoError(err)
		ts.EqualValues(want, list)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("New orders not found", func() {
		stOrder := stMocks.NewOrderer(ts.T())
		stOrder.EXPECT().GetByID(ctx, orderID).Return(models.Order{}, store.ErrNotFound)
		ts.setStorage(stOrder, nil, nil)

		_, err := ts.service.GetByID(ctx, orderID)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrOrderNotFound)
		ts.storage.AssertExpectations(ts.T())
	})
}
