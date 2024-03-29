// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	context "context"

	decimal "github.com/shopspring/decimal"
	mock "github.com/stretchr/testify/mock"

	models "github.com/fishus/go-advanced-gophermart/pkg/models"
)

// Loyaltier is an autogenerated mock type for the Loyaltier type
type Loyaltier struct {
	mock.Mock
}

type Loyaltier_Expecter struct {
	mock *mock.Mock
}

func (_m *Loyaltier) EXPECT() *Loyaltier_Expecter {
	return &Loyaltier_Expecter{mock: &_m.Mock}
}

// AddWithdraw provides a mock function with given fields: ctx, userID, orderNum, withdraw
func (_m *Loyaltier) AddWithdraw(ctx context.Context, userID models.UserID, orderNum string, withdraw decimal.Decimal) error {
	ret := _m.Called(ctx, userID, orderNum, withdraw)

	if len(ret) == 0 {
		panic("no return value specified for AddWithdraw")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID, string, decimal.Decimal) error); ok {
		r0 = rf(ctx, userID, orderNum, withdraw)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Loyaltier_AddWithdraw_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddWithdraw'
type Loyaltier_AddWithdraw_Call struct {
	*mock.Call
}

// AddWithdraw is a helper method to define mock.On call
//   - ctx context.Context
//   - userID models.UserID
//   - orderNum string
//   - withdraw decimal.Decimal
func (_e *Loyaltier_Expecter) AddWithdraw(ctx interface{}, userID interface{}, orderNum interface{}, withdraw interface{}) *Loyaltier_AddWithdraw_Call {
	return &Loyaltier_AddWithdraw_Call{Call: _e.mock.On("AddWithdraw", ctx, userID, orderNum, withdraw)}
}

func (_c *Loyaltier_AddWithdraw_Call) Run(run func(ctx context.Context, userID models.UserID, orderNum string, withdraw decimal.Decimal)) *Loyaltier_AddWithdraw_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.UserID), args[2].(string), args[3].(decimal.Decimal))
	})
	return _c
}

func (_c *Loyaltier_AddWithdraw_Call) Return(_a0 error) *Loyaltier_AddWithdraw_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Loyaltier_AddWithdraw_Call) RunAndReturn(run func(context.Context, models.UserID, string, decimal.Decimal) error) *Loyaltier_AddWithdraw_Call {
	_c.Call.Return(run)
	return _c
}

// UserBalance provides a mock function with given fields: _a0, _a1
func (_m *Loyaltier) UserBalance(_a0 context.Context, _a1 models.UserID) (models.LoyaltyBalance, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for UserBalance")
	}

	var r0 models.LoyaltyBalance
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID) (models.LoyaltyBalance, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID) models.LoyaltyBalance); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(models.LoyaltyBalance)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.UserID) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Loyaltier_UserBalance_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UserBalance'
type Loyaltier_UserBalance_Call struct {
	*mock.Call
}

// UserBalance is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 models.UserID
func (_e *Loyaltier_Expecter) UserBalance(_a0 interface{}, _a1 interface{}) *Loyaltier_UserBalance_Call {
	return &Loyaltier_UserBalance_Call{Call: _e.mock.On("UserBalance", _a0, _a1)}
}

func (_c *Loyaltier_UserBalance_Call) Run(run func(_a0 context.Context, _a1 models.UserID)) *Loyaltier_UserBalance_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.UserID))
	})
	return _c
}

func (_c *Loyaltier_UserBalance_Call) Return(_a0 models.LoyaltyBalance, _a1 error) *Loyaltier_UserBalance_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Loyaltier_UserBalance_Call) RunAndReturn(run func(context.Context, models.UserID) (models.LoyaltyBalance, error)) *Loyaltier_UserBalance_Call {
	_c.Call.Return(run)
	return _c
}

// UserWithdrawals provides a mock function with given fields: _a0, _a1
func (_m *Loyaltier) UserWithdrawals(_a0 context.Context, _a1 models.UserID) ([]models.LoyaltyHistory, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for UserWithdrawals")
	}

	var r0 []models.LoyaltyHistory
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID) ([]models.LoyaltyHistory, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID) []models.LoyaltyHistory); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.LoyaltyHistory)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.UserID) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Loyaltier_UserWithdrawals_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UserWithdrawals'
type Loyaltier_UserWithdrawals_Call struct {
	*mock.Call
}

// UserWithdrawals is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 models.UserID
func (_e *Loyaltier_Expecter) UserWithdrawals(_a0 interface{}, _a1 interface{}) *Loyaltier_UserWithdrawals_Call {
	return &Loyaltier_UserWithdrawals_Call{Call: _e.mock.On("UserWithdrawals", _a0, _a1)}
}

func (_c *Loyaltier_UserWithdrawals_Call) Run(run func(_a0 context.Context, _a1 models.UserID)) *Loyaltier_UserWithdrawals_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.UserID))
	})
	return _c
}

func (_c *Loyaltier_UserWithdrawals_Call) Return(_a0 []models.LoyaltyHistory, _a1 error) *Loyaltier_UserWithdrawals_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Loyaltier_UserWithdrawals_Call) RunAndReturn(run func(context.Context, models.UserID) ([]models.LoyaltyHistory, error)) *Loyaltier_UserWithdrawals_Call {
	_c.Call.Return(run)
	return _c
}

// NewLoyaltier creates a new instance of Loyaltier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLoyaltier(t interface {
	mock.TestingT
	Cleanup(func())
}) *Loyaltier {
	mock := &Loyaltier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
