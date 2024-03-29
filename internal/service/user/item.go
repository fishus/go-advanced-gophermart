package user

import (
	"context"
	"errors"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (s *service) GetByID(ctx context.Context, id models.UserID) (user models.User, err error) {
	user, err = s.storage.User().GetByID(ctx, id)
	if err != nil && errors.Is(err, store.ErrNotFound) {
		err = serviceErr.ErrUserNotFound
	}
	return
}
