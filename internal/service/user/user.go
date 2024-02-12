package user

import (
	"context"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (s *service) UserByID(ctx context.Context, id models.UserID) (*models.User, error) {
	user, err := s.storage.UserByID(ctx, id)
	if err != nil {
		return user, err
	}
	return user, nil
}
