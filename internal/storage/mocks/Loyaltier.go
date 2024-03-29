// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	context "context"

	decimal "github.com/shopspring/decimal"
	mock "github.com/stretchr/testify/mock"

	models "github.com/fishus/go-advanced-gophermart/pkg/models"

	pgx "github.com/jackc/pgx/v5"
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

// BalanceByUser provides a mock function with given fields: _a0, _a1
func (_m *Loyaltier) BalanceByUser(_a0 context.Context, _a1 models.UserID) (models.LoyaltyBalance, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for BalanceByUser")
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

// Loyaltier_BalanceByUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BalanceByUser'
type Loyaltier_BalanceByUser_Call struct {
	*mock.Call
}

// BalanceByUser is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 models.UserID
func (_e *Loyaltier_Expecter) BalanceByUser(_a0 interface{}, _a1 interface{}) *Loyaltier_BalanceByUser_Call {
	return &Loyaltier_BalanceByUser_Call{Call: _e.mock.On("BalanceByUser", _a0, _a1)}
}

func (_c *Loyaltier_BalanceByUser_Call) Run(run func(_a0 context.Context, _a1 models.UserID)) *Loyaltier_BalanceByUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.UserID))
	})
	return _c
}

func (_c *Loyaltier_BalanceByUser_Call) Return(_a0 models.LoyaltyBalance, _a1 error) *Loyaltier_BalanceByUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Loyaltier_BalanceByUser_Call) RunAndReturn(run func(context.Context, models.UserID) (models.LoyaltyBalance, error)) *Loyaltier_BalanceByUser_Call {
	_c.Call.Return(run)
	return _c
}

// BalanceUpdate provides a mock function with given fields: _a0, _a1, _a2
func (_m *Loyaltier) BalanceUpdate(_a0 context.Context, _a1 pgx.Tx, _a2 models.LoyaltyBalance) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for BalanceUpdate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx, models.LoyaltyBalance) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Loyaltier_BalanceUpdate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BalanceUpdate'
type Loyaltier_BalanceUpdate_Call struct {
	*mock.Call
}

// BalanceUpdate is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 pgx.Tx
//   - _a2 models.LoyaltyBalance
func (_e *Loyaltier_Expecter) BalanceUpdate(_a0 interface{}, _a1 interface{}, _a2 interface{}) *Loyaltier_BalanceUpdate_Call {
	return &Loyaltier_BalanceUpdate_Call{Call: _e.mock.On("BalanceUpdate", _a0, _a1, _a2)}
}

func (_c *Loyaltier_BalanceUpdate_Call) Run(run func(_a0 context.Context, _a1 pgx.Tx, _a2 models.LoyaltyBalance)) *Loyaltier_BalanceUpdate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(pgx.Tx), args[2].(models.LoyaltyBalance))
	})
	return _c
}

func (_c *Loyaltier_BalanceUpdate_Call) Return(_a0 error) *Loyaltier_BalanceUpdate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Loyaltier_BalanceUpdate_Call) RunAndReturn(run func(context.Context, pgx.Tx, models.LoyaltyBalance) error) *Loyaltier_BalanceUpdate_Call {
	_c.Call.Return(run)
	return _c
}

// HistoryAdd provides a mock function with given fields: _a0, _a1, _a2
func (_m *Loyaltier) HistoryAdd(_a0 context.Context, _a1 pgx.Tx, _a2 models.LoyaltyHistory) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for HistoryAdd")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx, models.LoyaltyHistory) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Loyaltier_HistoryAdd_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HistoryAdd'
type Loyaltier_HistoryAdd_Call struct {
	*mock.Call
}

// HistoryAdd is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 pgx.Tx
//   - _a2 models.LoyaltyHistory
func (_e *Loyaltier_Expecter) HistoryAdd(_a0 interface{}, _a1 interface{}, _a2 interface{}) *Loyaltier_HistoryAdd_Call {
	return &Loyaltier_HistoryAdd_Call{Call: _e.mock.On("HistoryAdd", _a0, _a1, _a2)}
}

func (_c *Loyaltier_HistoryAdd_Call) Run(run func(_a0 context.Context, _a1 pgx.Tx, _a2 models.LoyaltyHistory)) *Loyaltier_HistoryAdd_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(pgx.Tx), args[2].(models.LoyaltyHistory))
	})
	return _c
}

func (_c *Loyaltier_HistoryAdd_Call) Return(_a0 error) *Loyaltier_HistoryAdd_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Loyaltier_HistoryAdd_Call) RunAndReturn(run func(context.Context, pgx.Tx, models.LoyaltyHistory) error) *Loyaltier_HistoryAdd_Call {
	_c.Call.Return(run)
	return _c
}

// HistoryByUser provides a mock function with given fields: _a0, _a1
func (_m *Loyaltier) HistoryByUser(_a0 context.Context, _a1 models.UserID) ([]models.LoyaltyHistory, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for HistoryByUser")
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

// Loyaltier_HistoryByUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HistoryByUser'
type Loyaltier_HistoryByUser_Call struct {
	*mock.Call
}

// HistoryByUser is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 models.UserID
func (_e *Loyaltier_Expecter) HistoryByUser(_a0 interface{}, _a1 interface{}) *Loyaltier_HistoryByUser_Call {
	return &Loyaltier_HistoryByUser_Call{Call: _e.mock.On("HistoryByUser", _a0, _a1)}
}

func (_c *Loyaltier_HistoryByUser_Call) Run(run func(_a0 context.Context, _a1 models.UserID)) *Loyaltier_HistoryByUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.UserID))
	})
	return _c
}

func (_c *Loyaltier_HistoryByUser_Call) Return(_a0 []models.LoyaltyHistory, _a1 error) *Loyaltier_HistoryByUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Loyaltier_HistoryByUser_Call) RunAndReturn(run func(context.Context, models.UserID) ([]models.LoyaltyHistory, error)) *Loyaltier_HistoryByUser_Call {
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
