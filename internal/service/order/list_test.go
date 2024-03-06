package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (ts *OrderServiceTestSuite) TestListNew() {
	ctx := context.Background()

	ts.Run("Positive case", func() {
		want := []models.Order{
			{
				ID:         models.OrderID(uuid.New().String()),
				UserID:     models.UserID(uuid.New().String()),
				Num:        "9305514466",
				Accrual:    decimal.NewFromFloat(0),
				Status:     models.OrderStatusNew,
				UploadedAt: time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
			},
			{
				ID:         models.OrderID(uuid.New().String()),
				UserID:     models.UserID(uuid.New().String()),
				Num:        "1206405415",
				Accrual:    decimal.NewFromFloat(0),
				Status:     models.OrderStatusNew,
				UploadedAt: time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
			},
		}
		mockCall := ts.storage.EXPECT().OrdersByFilter(ctx, 0, mock.Anything, mock.Anything).Return(want, nil)
		defer mockCall.Unset()
		list, err := ts.service.ListNew(ctx)
		ts.NoError(err)
		ts.EqualValues(want, list)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("New orders not found", func() {
		mockCall := ts.storage.EXPECT().OrdersByFilter(ctx, 0, mock.Anything, mock.Anything).Return([]models.Order{}, store.ErrNotFound)
		defer mockCall.Unset()
		list, err := ts.service.ListNew(ctx)
		ts.NoError(err)
		ts.Equal(0, len(list))
		ts.storage.AssertExpectations(ts.T())
	})
}

func (ts *OrderServiceTestSuite) TestListProcessing() {
	ctx := context.Background()

	limit := 10

	ts.Run("Positive case", func() {
		want := []models.Order{
			{
				ID:         models.OrderID(uuid.New().String()),
				UserID:     models.UserID(uuid.New().String()),
				Num:        "9305514466",
				Accrual:    decimal.NewFromFloat(0),
				Status:     models.OrderStatusProcessing,
				UploadedAt: time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
			},
			{
				ID:         models.OrderID(uuid.New().String()),
				UserID:     models.UserID(uuid.New().String()),
				Num:        "1206405415",
				Accrual:    decimal.NewFromFloat(0),
				Status:     models.OrderStatusProcessing,
				UploadedAt: time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
			},
		}
		mockCall := ts.storage.EXPECT().OrdersByFilter(ctx, limit, mock.Anything, mock.Anything).Return(want, nil)
		defer mockCall.Unset()
		list, err := ts.service.ListProcessing(ctx, limit)
		ts.NoError(err)
		ts.EqualValues(want, list)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("Orders in processing not found", func() {
		mockCall := ts.storage.EXPECT().OrdersByFilter(ctx, limit, mock.Anything, mock.Anything).Return([]models.Order{}, store.ErrNotFound)
		defer mockCall.Unset()
		list, err := ts.service.ListProcessing(ctx, limit)
		ts.NoError(err)
		ts.Equal(0, len(list))
		ts.storage.AssertExpectations(ts.T())
	})
}

func (ts *OrderServiceTestSuite) TestListByUser() {
	ctx := context.Background()

	userID := models.UserID(uuid.New().String())

	ts.Run("Positive case", func() {
		want := []models.Order{
			{
				ID:         models.OrderID(uuid.New().String()),
				UserID:     userID,
				Num:        "9305514466",
				Accrual:    decimal.NewFromFloat(0),
				Status:     models.OrderStatusNew,
				UploadedAt: time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
			},
			{
				ID:         models.OrderID(uuid.New().String()),
				UserID:     userID,
				Num:        "1206405415",
				Accrual:    decimal.NewFromFloat(0),
				Status:     models.OrderStatusProcessing,
				UploadedAt: time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
			},
			{
				ID:         models.OrderID(uuid.New().String()),
				UserID:     userID,
				Num:        "1853241857",
				Accrual:    decimal.NewFromFloatWithExponent(123.456, -5),
				Status:     models.OrderStatusProcessed,
				UploadedAt: time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
			},
		}
		mockCall := ts.storage.EXPECT().OrdersByFilter(ctx, 0, mock.Anything, mock.Anything).Return(want, nil)
		defer mockCall.Unset()
		list, err := ts.service.ListByUser(ctx, userID)
		ts.NoError(err)
		ts.EqualValues(want, list)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("New orders not found", func() {
		mockCall := ts.storage.EXPECT().OrdersByFilter(ctx, 0, mock.Anything, mock.Anything).Return([]models.Order{}, store.ErrNotFound)
		defer mockCall.Unset()
		list, err := ts.service.ListByUser(ctx, userID)
		ts.NoError(err)
		ts.Equal(0, len(list))
		ts.storage.AssertExpectations(ts.T())
	})
}
