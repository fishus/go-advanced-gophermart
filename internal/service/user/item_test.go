package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	stMocks "github.com/fishus/go-advanced-gophermart/internal/storage/mocks"
)

func (ts *UserServiceTestSuite) TestGetByID() {
	ctx := context.Background()

	ts.Run("Return user by ID", func() {
		userID := models.UserID(uuid.New().String())
		username := make([]byte, 10)
		_, err := rand.Read(username)
		ts.Require().NoError(err)
		want := models.User{
			ID:        userID,
			Username:  hex.EncodeToString(username),
			Password:  hex.EncodeToString(username),
			CreatedAt: time.Now().UTC().Round(time.Second),
		}

		stUser := stMocks.NewUserer(ts.T())
		stUser.EXPECT().GetByID(ctx, userID).Return(want, nil)
		ts.setStorage(nil, stUser, nil)

		user, err := ts.service.GetByID(ctx, userID)
		ts.NoError(err)
		ts.EqualValues(want, user)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("User not found", func() {
		userID := models.UserID(uuid.New().String())
		want := models.User{}

		stUser := stMocks.NewUserer(ts.T())
		stUser.EXPECT().GetByID(ctx, userID).Return(want, store.ErrNotFound)
		ts.setStorage(nil, stUser, nil)

		_, err := ts.service.GetByID(ctx, userID)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrUserNotFound)
		ts.storage.AssertExpectations(ts.T())
	})
}
