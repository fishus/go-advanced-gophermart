package user

import (
	"context"
	"errors"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

// Register Регистрация пользователя
func (s *service) Register(ctx context.Context, user models.User) (userID models.UserID, err error) {
	err = validateUser(user)
	if err != nil {
		return
	}

	userID, err = s.storage.User().Add(ctx, user)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			err = serviceErr.ErrUserAlreadyExists
			return
		}
		if errors.Is(err, store.ErrIncorrectData) {
			err = serviceErr.ErrIncorrectData
			return
		}
	}
	return
}
