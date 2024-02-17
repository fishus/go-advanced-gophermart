package user

import (
	"context"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (s *service) UserByID(ctx context.Context, id models.UserID) (user models.User, err error) {
	user, err = s.storage.UserByID(ctx, id)
	return
}
