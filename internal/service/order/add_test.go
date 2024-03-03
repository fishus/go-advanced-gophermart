package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
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
		mockCall := ts.storage.EXPECT().OrderAdd(ctx, data).Return(wantID, nil)
		defer mockCall.Unset()

		id, err := ts.service.Add(ctx, userID, data.Num)
		ts.NoError(err)
		ts.Equal(wantID, id)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("Order already exists (my own)", func() {
		mockCall1 := ts.storage.EXPECT().OrderAdd(ctx, data).Return("", store.ErrAlreadyExists)
		defer mockCall1.Unset()

		mockCall2 := ts.storage.EXPECT().OrderByFilter(ctx, mock.Anything).Return(data, nil)
		defer mockCall2.Unset()

		_, err := ts.service.Add(ctx, userID, data.Num)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrOrderAlreadyExists)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("Order already exists (non-owned)", func() {
		dataExists := data
		dataExists.UserID = models.UserID(uuid.New().String())
		mockCall1 := ts.storage.EXPECT().OrderAdd(ctx, data).Return("", store.ErrAlreadyExists)
		defer mockCall1.Unset()

		// mock.FunctionalOptions(store.WithOrderNum(data.Num))
		mockCall2 := ts.storage.EXPECT().OrderByFilter(ctx, mock.Anything).Return(dataExists, nil)
		defer mockCall2.Unset()

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
