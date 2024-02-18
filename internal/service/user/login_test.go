package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (ts *UserServiceTestSuite) TestLogin() {
	ctx := context.Background()

	storage := new(store.MockStorage)

	wantID := models.UserID(uuid.New().String())
	username := make([]byte, 10)
	_, err := rand.Read(username)
	ts.Require().NoError(err)
	data := models.User{
		Username: hex.EncodeToString(username),
		Password: hex.EncodeToString(username),
	}

	ts.Run("Positive case", func() {
		mockCall := storage.On("UserLogin", ctx, data).Return(wantID, nil)
		service := New(&Config{}, storage)

		id, err := service.Login(ctx, data)
		ts.NoError(err)
		ts.Equal(wantID, id)
		storage.AssertExpectations(ts.T())
		mockCall.Unset()
	})

	ts.Run("User not found", func() {
		data := models.User{
			Username: "test",
			Password: "test",
		}
		mockCall := storage.On("UserLogin", ctx, data).Return(models.UserID(""), store.ErrNotFound)
		service := New(&Config{}, storage)

		_, err := service.Login(ctx, data)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrUserNotFound)
		storage.AssertExpectations(ts.T())
		mockCall.Unset()
	})

	ts.Run("Not valid", func() {
		data := models.User{}
		service := New(&Config{}, storage)

		_, err := service.Login(ctx, data)
		ts.Error(err)
		var ve *serviceErr.ValidationError
		ts.ErrorAs(err, &ve)
		storage.AssertNotCalled(ts.T(), "UserLogin")
	})
}
