package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	stMocks "github.com/fishus/go-advanced-gophermart/internal/storage/mocks"
)

func (ts *OrderServiceTestSuite) TestAdd() {
	ctx := context.Background()

	userID := models.UserID(uuid.New().String())
	wantID := models.OrderID(uuid.New().String())
	data := models.Order{
		UserID: userID,
		Num:    "2313224962",
		Status: models.OrderStatusNew,
	}

	ts.Run("Positive case", func() {
		stOrder := stMocks.NewOrderer(ts.T())
		stOrder.EXPECT().Add(ctx, data).Return(wantID, nil)
		ts.setStorage(stOrder, nil, nil)

		id, err := ts.service.Add(ctx, userID, data.Num)
		ts.NoError(err)
		ts.Equal(wantID, id)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("Order already exists (my own)", func() {
		stOrder := stMocks.NewOrderer(ts.T())

		stOrder.EXPECT().Add(ctx, data).Return("", store.ErrAlreadyExists)
		stOrder.EXPECT().GetByFilter(ctx, mock.Anything).Return(data, nil)

		ts.setStorage(stOrder, nil, nil)

		_, err := ts.service.Add(ctx, userID, data.Num)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrOrderAlreadyExists)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("Order already exists (non-owned)", func() {
		stOrder := stMocks.NewOrderer(ts.T())

		dataExists := data
		dataExists.UserID = models.UserID(uuid.New().String())
		stOrder.EXPECT().Add(ctx, data).Return("", store.ErrAlreadyExists)

		// mock.FunctionalOptions(store.WithOrderNum(data.Num))
		stOrder.EXPECT().GetByFilter(ctx, mock.Anything).Return(dataExists, nil)

		ts.setStorage(stOrder, nil, nil)

		_, err := ts.service.Add(ctx, userID, data.Num)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrOrderWrongOwner)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("Not valid", func() {
		_, err := ts.service.Add(ctx, "", "")
		ts.Error(err)
		var ve *serviceErr.ValidationError
		ts.ErrorAs(err, &ve)
	})
}
