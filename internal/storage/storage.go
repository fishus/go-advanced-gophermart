package storage

import (
	"context"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=UserStorager  --with-expecter
type UserStorager interface {
	UserAdd(context.Context, models.User) (models.UserID, error)
	UserLogin(context.Context, models.User) (models.UserID, error)
	UserByID(context.Context, models.UserID) (models.User, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=OrderStorager  --with-expecter
type OrderStorager interface {
	OrderAdd(context.Context, models.Order) (models.OrderID, error)
	OrderByID(context.Context, models.OrderID) (models.Order, error)
	OrderByFilter(context.Context, ...OrderFilter) (models.Order, error)
	OrdersByFilter(ctx context.Context, limit int, filters ...OrderFilter) ([]models.Order, error)
	OrderUpdateStatus(context.Context, models.OrderID, models.OrderStatus) error
	OrderAddAccrual(ctx context.Context, orderID models.OrderID, accrual float64) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=LoyaltyBalancer  --with-expecter
type LoyaltyBalancer interface {
	LoyaltyBalanceByUser(context.Context, models.UserID) (models.LoyaltyBalance, error)
	LoyaltyAddWithdraw(ctx context.Context, userID models.UserID, orderNum string, withdraw float64) error
	LoyaltyHistoryByUser(context.Context, models.UserID) ([]models.LoyaltyHistory, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=Storager  --with-expecter
type Storager interface {
	UserStorager
	OrderStorager
	LoyaltyBalancer
}
