package user

import (
	"context"
	"github.com/gookit/validate"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (s *service) Login(ctx context.Context, user models.User) (models.User, error) {
	v := validate.Struct(user)
	if !v.Validate() {
		return user, serviceErr.NewValidationError(v.Errors)
	}

	user, err := s.storage.UserLogin(ctx, user)
	if err != nil {
		return user, err
	}
	return user, nil
}
