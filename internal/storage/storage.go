package storage

import (
	"context"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

type UserStorager interface {
	UserAdd(context.Context, models.User) (models.UserID, error)
	UserLogin(context.Context, models.User) (models.UserID, error)
	UserByID(context.Context, models.UserID) (models.User, error)
}

type OrderStorager interface {
	OrderAdd(context.Context, models.Order) (models.OrderID, error)
	OrderByID(context.Context, models.OrderID) (models.Order, error)
	OrderByFilter(context.Context, ...OrderFilter) (models.Order, error)
	OrdersByFilter(ctx context.Context, limit int, filters ...OrderFilter) ([]models.Order, error)
	OrderUpdateStatus(context.Context, models.OrderID, models.OrderStatus) error
	OrderAddAccrual(ctx context.Context, order models.Order, accrual float64) error
}

type Storager interface {
	UserStorager
	OrderStorager
}
