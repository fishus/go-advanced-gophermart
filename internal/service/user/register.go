package user

import (
	"context"
	"errors"

	"github.com/gookit/validate"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

// Register Регистрация пользователя
func (s *service) Register(ctx context.Context, user models.User) (models.UserID, error) {
	var userID models.UserID

	v := validate.Struct(user)
	if !v.Validate() {
		return userID, serviceErr.NewValidationError(v.Errors)
	}

	userID, err := s.storage.UserAdd(ctx, user)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			err = serviceErr.ErrUserAlreadyExists
		}
		return userID, err
	}

	return userID, nil
}
