package user

import (
	"context"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
	"github.com/gookit/validate"
)

func (s *service) Register(ctx context.Context, user models.User) (models.User, error) {
	v := validate.Struct(user)
	if !v.Validate() {
		return user, serviceErr.NewValidationError(v.Errors)
	}

	user, err := s.storage.UserAdd(ctx, user)
	if err != nil {
		return user, err
	}
	return user, nil
}
