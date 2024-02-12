package user

import (
	"context"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
	"github.com/gookit/validate"
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
		return userID, err
	}

	return userID, nil
}
