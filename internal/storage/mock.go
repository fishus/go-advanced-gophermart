package storage

import (
	"context"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

var _ Storager = (*MockStorage)(nil)

func (m *MockStorage) UserAdd(ctx context.Context, user models.User) (models.UserID, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(models.UserID), args.Error(1)
}
func (m *MockStorage) UserLogin(ctx context.Context, user models.User) (models.UserID, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(models.UserID), args.Error(1)
}
func (m *MockStorage) UserByID(ctx context.Context, id models.UserID) (models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.User), args.Error(1)
}
func (m *MockStorage) OrderAdd(ctx context.Context, order models.Order) (models.OrderID, error) {
	args := m.Called(ctx, order)
	return args.Get(0).(models.OrderID), args.Error(1)
}
func (m *MockStorage) OrderByID(ctx context.Context, id models.OrderID) (models.Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Order), args.Error(1)
}
func (m *MockStorage) OrderByFilter(ctx context.Context, filters ...OrderFilter) (models.Order, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(models.Order), args.Error(1)
}
func (m *MockStorage) OrdersByFilter(ctx context.Context, limit int, filters ...OrderFilter) ([]models.Order, error) {
	args := m.Called(ctx, limit, filters)
	return args.Get(0).([]models.Order), args.Error(1)
}
func (m *MockStorage) OrderUpdateStatus(ctx context.Context, id models.OrderID, status models.OrderStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}
func (m *MockStorage) OrderAddAccrual(ctx context.Context, order models.Order, accrual float64) error {
	args := m.Called(ctx, order, accrual)
	return args.Error(0)
}
func (m *MockStorage) LoyaltyBalanceByUser(ctx context.Context, userID models.UserID) (models.LoyaltyBalance, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(models.LoyaltyBalance), args.Error(1)
}
func (m *MockStorage) LoyaltyAddWithdraw(ctx context.Context, userID models.UserID, orderNum string, withdraw float64) error {
	args := m.Called(ctx, userID, orderNum, withdraw)
	return args.Error(0)
}
