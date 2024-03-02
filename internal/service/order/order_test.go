package order

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (ts *OrderServiceTestSuite) TestOrderByID() {
	ctx := context.Background()

	orderID := models.OrderID(uuid.New().String())

	ts.Run("Positive case", func() {
		want := models.Order{
			ID:         orderID,
			UserID:     models.UserID(uuid.New().String()),
			Num:        "9400781309",
			Accrual:    0,
			Status:     models.OrderStatusNew,
			UploadedAt: time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		}
		mockCall := ts.storage.On("OrderByID", ctx, orderID).Return(want, nil)
		list, err := ts.service.OrderByID(ctx, orderID)
		ts.NoError(err)
		ts.EqualValues(want, list)
		ts.storage.AssertExpectations(ts.T())
		mockCall.Unset()
	})

	ts.Run("New orders not found", func() {
		mockCall := ts.storage.On("OrderByID", ctx, orderID).Return(models.Order{}, store.ErrNotFound)
		_, err := ts.service.OrderByID(ctx, orderID)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrOrderNotFound)
		ts.storage.AssertExpectations(ts.T())
		mockCall.Unset()
	})
}
