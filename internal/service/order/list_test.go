package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/app/config"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	stMocks "github.com/fishus/go-advanced-gophermart/internal/storage/mocks"
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

		stOrder := stMocks.NewOrderer(ts.T())
		stOrder.EXPECT().ListByFilter(ctx, 0, mock.MatchedBy(func(filter store.OrderFilter) bool {
			f := &store.OrderFilters{}
			filter(f)
			return len(f.Statuses) == 1 && f.Statuses[0] == models.OrderStatusNew
		}), mock.MatchedBy(func(filter store.OrderFilter) bool {
			f := &store.OrderFilters{}
			filter(f)
			return len(f.OrderBy) == 1 && f.OrderBy[0].Field == store.OrderByUploadedAt && f.OrderBy[0].Dir == store.OrderByAsc
		})).Return(want, nil)
		ts.setStorage(stOrder, nil, nil)

		list, err := ts.service.ListNew(ctx)
		ts.NoError(err)
		ts.EqualValues(want, list)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("New orders not found", func() {
		stOrder := stMocks.NewOrderer(ts.T())
		stOrder.EXPECT().ListByFilter(ctx, 0, mock.MatchedBy(func(filter store.OrderFilter) bool {
			f := &store.OrderFilters{}
			filter(f)
			return len(f.Statuses) == 1 && f.Statuses[0] == models.OrderStatusNew
		}), mock.MatchedBy(func(filter store.OrderFilter) bool {
			f := &store.OrderFilters{}
			filter(f)
			return len(f.OrderBy) == 1 && f.OrderBy[0].Field == store.OrderByUploadedAt && f.OrderBy[0].Dir == store.OrderByAsc
		})).Return([]models.Order{}, store.ErrNotFound)
		ts.setStorage(stOrder, nil, nil)

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

		stOrder := stMocks.NewOrderer(ts.T())
		stOrder.EXPECT().ListByFilter(ctx, limit, mock.MatchedBy(func(filter store.OrderFilter) bool {
			f := &store.OrderFilters{}
			filter(f)
			return len(f.Statuses) == 1 && f.Statuses[0] == models.OrderStatusProcessing
		}), mock.MatchedBy(func(filter store.OrderFilter) bool {
			f := &store.OrderFilters{}
			filter(f)
			return len(f.OrderBy) == 1 && f.OrderBy[0].Field == store.OrderByUpdatedAt && f.OrderBy[0].Dir == store.OrderByAsc
		})).Return(want, nil)
		ts.setStorage(stOrder, nil, nil)

		list, err := ts.service.ListProcessing(ctx, limit)
		ts.NoError(err)
		ts.EqualValues(want, list)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("Orders in processing not found", func() {
		stOrder := stMocks.NewOrderer(ts.T())
		stOrder.EXPECT().ListByFilter(ctx, limit, mock.MatchedBy(func(filter store.OrderFilter) bool {
			f := &store.OrderFilters{}
			filter(f)
			return len(f.Statuses) == 1 && f.Statuses[0] == models.OrderStatusProcessing
		}), mock.MatchedBy(func(filter store.OrderFilter) bool {
			f := &store.OrderFilters{}
			filter(f)
			return len(f.OrderBy) == 1 && f.OrderBy[0].Field == store.OrderByUpdatedAt && f.OrderBy[0].Dir == store.OrderByAsc
		})).Return([]models.Order{}, store.ErrNotFound)
		ts.setStorage(stOrder, nil, nil)

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
				Accrual:    decimal.NewFromFloatWithExponent(123.456, -config.DecimalExponent),
				Status:     models.OrderStatusProcessed,
				UploadedAt: time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
			},
		}

		stOrder := stMocks.NewOrderer(ts.T())
		stOrder.EXPECT().ListByFilter(ctx, 0, mock.MatchedBy(func(filter store.OrderFilter) bool {
			f := &store.OrderFilters{}
			filter(f)
			return f.UserID == userID
		}), mock.MatchedBy(func(filter store.OrderFilter) bool {
			f := &store.OrderFilters{}
			filter(f)
			return len(f.OrderBy) == 1 && f.OrderBy[0].Field == store.OrderByUploadedAt && f.OrderBy[0].Dir == store.OrderByAsc
		})).Return(want, nil)
		ts.setStorage(stOrder, nil, nil)

		list, err := ts.service.ListByUser(ctx, userID)
		ts.NoError(err)
		ts.EqualValues(want, list)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("New orders not found", func() {
		stOrder := stMocks.NewOrderer(ts.T())
		stOrder.EXPECT().ListByFilter(ctx, 0, mock.MatchedBy(func(filter store.OrderFilter) bool {
			f := &store.OrderFilters{}
			filter(f)
			return f.UserID == userID
		}), mock.MatchedBy(func(filter store.OrderFilter) bool {
			f := &store.OrderFilters{}
			filter(f)
			return len(f.OrderBy) == 1 && f.OrderBy[0].Field == store.OrderByUploadedAt && f.OrderBy[0].Dir == store.OrderByAsc
		})).Return([]models.Order{}, store.ErrNotFound)
		ts.setStorage(stOrder, nil, nil)

		list, err := ts.service.ListByUser(ctx, userID)
		ts.NoError(err)
		ts.Equal(0, len(list))
		ts.storage.AssertExpectations(ts.T())
	})
}
