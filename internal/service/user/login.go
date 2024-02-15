package user

import (
	"context"
	"errors"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

// Login Аутентификация пользователя
func (s *service) Login(ctx context.Context, user models.User) (userID models.UserID, err error) {
	err = validateUser(user)
	if err != nil {
		return
	}

	userID, err = s.storage.UserLogin(ctx, user)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			err = serviceErr.ErrUserAlreadyExists
		}
	}
	return
}
