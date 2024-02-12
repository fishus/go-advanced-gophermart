package storage

import (
	"context"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

type UserStorager interface {
	UserAdd(context.Context, models.User) (models.User, error)
	UserLogin(context.Context, models.User) (models.User, error)
}

type Storager interface {
	UserStorager
}
