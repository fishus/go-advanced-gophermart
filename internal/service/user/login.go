package user

import (
	"context"
	
	"github.com/gookit/validate"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (s *service) Login(ctx context.Context, user models.User) (models.UserID, error) {
	var userID models.UserID

	v := validate.Struct(user)
	if !v.Validate() {
		return userID, serviceErr.NewValidationError(v.Errors)
	}

	userID, err := s.storage.UserLogin(ctx, user)
	if err != nil {
		return userID, err
	}

	return userID, nil
}
