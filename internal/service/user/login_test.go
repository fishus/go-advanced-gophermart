package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	stMocks "github.com/fishus/go-advanced-gophermart/internal/storage/mocks"
)

func (ts *UserServiceTestSuite) TestLogin() {
	ctx := context.Background()

	wantID := models.UserID(uuid.New().String())
	username := make([]byte, 10)
	_, err := rand.Read(username)
	ts.Require().NoError(err)
	data := models.User{
		Username: hex.EncodeToString(username),
		Password: hex.EncodeToString(username),
	}

	ts.Run("Positive case", func() {
		stUser := stMocks.NewUserer(ts.T())
		stUser.EXPECT().Login(ctx, data).Return(wantID, nil)
		ts.setStorage(nil, stUser, nil)

		id, err := ts.service.Login(ctx, data)
		ts.NoError(err)
		ts.Equal(wantID, id)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("User not found", func() {
		data := models.User{
			Username: "test",
			Password: "test",
		}

		stUser := stMocks.NewUserer(ts.T())
		stUser.EXPECT().Login(ctx, data).Return("", store.ErrNotFound)
		ts.setStorage(nil, stUser, nil)

		_, err := ts.service.Login(ctx, data)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrUserNotFound)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("Not valid", func() {
		data := models.User{}

		_, err := ts.service.Login(ctx, data)
		ts.Error(err)
		var ve *serviceErr.ValidationError
		ts.ErrorAs(err, &ve)
	})
}
