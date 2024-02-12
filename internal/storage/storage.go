package storage

import (
	"context"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

type UserStorager interface {
	UserAdd(context.Context, models.User) (models.UserID, error)
	UserLogin(context.Context, models.User) (models.UserID, error)
	UserByID(context.Context, models.UserID) (*models.User, error)
}

type Storager interface {
	UserStorager
}
