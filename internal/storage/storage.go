package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=Userer  --with-expecter
type Userer interface {
	GetByID(context.Context, models.UserID) (models.User, error)
	Add(context.Context, models.User) (models.UserID, error)
	Login(context.Context, models.User) (models.UserID, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=Orderer  --with-expecter
type Orderer interface {
	GetByID(context.Context, models.OrderID) (models.Order, error)
	GetByFilter(context.Context, ...OrderFilter) (models.Order, error)
	Add(context.Context, models.Order) (models.OrderID, error)
	UpdateStatus(context.Context, models.OrderID, models.OrderStatus) error
	AddAccrual(ctx context.Context, orderID models.OrderID, accrual decimal.Decimal) error
	ListByFilter(ctx context.Context, limit int, filters ...OrderFilter) ([]models.Order, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=Loyaltier  --with-expecter
type Loyaltier interface {
	BalanceUpdate(context.Context, pgx.Tx, models.LoyaltyBalance) error
	AddWithdraw(ctx context.Context, userID models.UserID, orderNum string, withdraw decimal.Decimal) error
	BalanceByUser(context.Context, models.UserID) (models.LoyaltyBalance, error)
	HistoryAdd(context.Context, pgx.Tx, models.LoyaltyHistory) error
	HistoryByUser(context.Context, models.UserID) ([]models.LoyaltyHistory, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=Storager  --with-expecter
type Storager interface {
	User() Userer
	Order() Orderer
	Loyalty() Loyaltier
}
