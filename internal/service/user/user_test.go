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
)

func (ts *UserServiceTestSuite) TestUserByID() {
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
		mockCall := ts.storage.EXPECT().UserByID(ctx, userID).Return(want, nil)
		defer mockCall.Unset()

		user, err := ts.service.UserByID(ctx, userID)
		ts.NoError(err)
		ts.EqualValues(want, user)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("User not found", func() {
		userID := models.UserID(uuid.New().String())
		want := models.User{}
		mockCall := ts.storage.EXPECT().UserByID(ctx, userID).Return(want, store.ErrNotFound)
		defer mockCall.Unset()

		_, err := ts.service.UserByID(ctx, userID)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrUserNotFound)
		ts.storage.AssertExpectations(ts.T())
	})
}
