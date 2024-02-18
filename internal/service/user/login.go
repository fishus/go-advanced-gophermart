package user

import (
	"context"
	"errors"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

// Login Аутентификация пользователя
func (s *service) Login(ctx context.Context, user models.User) (userID models.UserID, err error) {
	err = validateUser(user)
	if err != nil {
		return
	}

	userID, err = s.storage.UserLogin(ctx, user)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			err = serviceErr.ErrUserNotFound
		}
	}
	return
}
