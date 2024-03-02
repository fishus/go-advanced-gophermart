// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/fishus/go-advanced-gophermart/pkg/models"
	mock "github.com/stretchr/testify/mock"
)

// Orderer is an autogenerated mock type for the Orderer type
type Orderer struct {
	mock.Mock
}

type Orderer_Expecter struct {
	mock *mock.Mock
}

func (_m *Orderer) EXPECT() *Orderer_Expecter {
	return &Orderer_Expecter{mock: &_m.Mock}
}

// Add provides a mock function with given fields: ctx, userID, orderNum
func (_m *Orderer) Add(ctx context.Context, userID models.UserID, orderNum string) (models.OrderID, error) {
	ret := _m.Called(ctx, userID, orderNum)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 models.OrderID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID, string) (models.OrderID, error)); ok {
		return rf(ctx, userID, orderNum)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID, string) models.OrderID); ok {
		r0 = rf(ctx, userID, orderNum)
	} else {
		r0 = ret.Get(0).(models.OrderID)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.UserID, string) error); ok {
		r1 = rf(ctx, userID, orderNum)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Orderer_Add_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Add'
type Orderer_Add_Call struct {
	*mock.Call
}

// Add is a helper method to define mock.On call
//   - ctx context.Context
//   - userID models.UserID
//   - orderNum string
func (_e *Orderer_Expecter) Add(ctx interface{}, userID interface{}, orderNum interface{}) *Orderer_Add_Call {
	return &Orderer_Add_Call{Call: _e.mock.On("Add", ctx, userID, orderNum)}
}

func (_c *Orderer_Add_Call) Run(run func(ctx context.Context, userID models.UserID, orderNum string)) *Orderer_Add_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.UserID), args[2].(string))
	})
	return _c
}

func (_c *Orderer_Add_Call) Return(_a0 models.OrderID, _a1 error) *Orderer_Add_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Orderer_Add_Call) RunAndReturn(run func(context.Context, models.UserID, string) (models.OrderID, error)) *Orderer_Add_Call {
	_c.Call.Return(run)
	return _c
}

// AddAccrual provides a mock function with given fields: ctx, id, accrual
func (_m *Orderer) AddAccrual(ctx context.Context, id models.OrderID, accrual float64) error {
	ret := _m.Called(ctx, id, accrual)

	if len(ret) == 0 {
		panic("no return value specified for AddAccrual")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.OrderID, float64) error); ok {
		r0 = rf(ctx, id, accrual)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Orderer_AddAccrual_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddAccrual'
type Orderer_AddAccrual_Call struct {
	*mock.Call
}

// AddAccrual is a helper method to define mock.On call
//   - ctx context.Context
//   - id models.OrderID
//   - accrual float64
func (_e *Orderer_Expecter) AddAccrual(ctx interface{}, id interface{}, accrual interface{}) *Orderer_AddAccrual_Call {
	return &Orderer_AddAccrual_Call{Call: _e.mock.On("AddAccrual", ctx, id, accrual)}
}

func (_c *Orderer_AddAccrual_Call) Run(run func(ctx context.Context, id models.OrderID, accrual float64)) *Orderer_AddAccrual_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.OrderID), args[2].(float64))
	})
	return _c
}

func (_c *Orderer_AddAccrual_Call) Return(_a0 error) *Orderer_AddAccrual_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Orderer_AddAccrual_Call) RunAndReturn(run func(context.Context, models.OrderID, float64) error) *Orderer_AddAccrual_Call {
	_c.Call.Return(run)
	return _c
}

// ListByUser provides a mock function with given fields: _a0, _a1
func (_m *Orderer) ListByUser(_a0 context.Context, _a1 models.UserID) ([]models.Order, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ListByUser")
	}

	var r0 []models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID) ([]models.Order, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID) []models.Order); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.UserID) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Orderer_ListByUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListByUser'
type Orderer_ListByUser_Call struct {
	*mock.Call
}

// ListByUser is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 models.UserID
func (_e *Orderer_Expecter) ListByUser(_a0 interface{}, _a1 interface{}) *Orderer_ListByUser_Call {
	return &Orderer_ListByUser_Call{Call: _e.mock.On("ListByUser", _a0, _a1)}
}

func (_c *Orderer_ListByUser_Call) Run(run func(_a0 context.Context, _a1 models.UserID)) *Orderer_ListByUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.UserID))
	})
	return _c
}

func (_c *Orderer_ListByUser_Call) Return(_a0 []models.Order, _a1 error) *Orderer_ListByUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Orderer_ListByUser_Call) RunAndReturn(run func(context.Context, models.UserID) ([]models.Order, error)) *Orderer_ListByUser_Call {
	_c.Call.Return(run)
	return _c
}

// ListNew provides a mock function with given fields: _a0
func (_m *Orderer) ListNew(_a0 context.Context) ([]models.Order, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for ListNew")
	}

	var r0 []models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]models.Order, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []models.Order); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Orderer_ListNew_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListNew'
type Orderer_ListNew_Call struct {
	*mock.Call
}

// ListNew is a helper method to define mock.On call
//   - _a0 context.Context
func (_e *Orderer_Expecter) ListNew(_a0 interface{}) *Orderer_ListNew_Call {
	return &Orderer_ListNew_Call{Call: _e.mock.On("ListNew", _a0)}
}

func (_c *Orderer_ListNew_Call) Run(run func(_a0 context.Context)) *Orderer_ListNew_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Orderer_ListNew_Call) Return(_a0 []models.Order, _a1 error) *Orderer_ListNew_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Orderer_ListNew_Call) RunAndReturn(run func(context.Context) ([]models.Order, error)) *Orderer_ListNew_Call {
	_c.Call.Return(run)
	return _c
}

// ListProcessing provides a mock function with given fields: ctx, limit
func (_m *Orderer) ListProcessing(ctx context.Context, limit int) ([]models.Order, error) {
	ret := _m.Called(ctx, limit)

	if len(ret) == 0 {
		panic("no return value specified for ListProcessing")
	}

	var r0 []models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) ([]models.Order, error)); ok {
		return rf(ctx, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) []models.Order); ok {
		r0 = rf(ctx, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Orderer_ListProcessing_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListProcessing'
type Orderer_ListProcessing_Call struct {
	*mock.Call
}

// ListProcessing is a helper method to define mock.On call
//   - ctx context.Context
//   - limit int
func (_e *Orderer_Expecter) ListProcessing(ctx interface{}, limit interface{}) *Orderer_ListProcessing_Call {
	return &Orderer_ListProcessing_Call{Call: _e.mock.On("ListProcessing", ctx, limit)}
}

func (_c *Orderer_ListProcessing_Call) Run(run func(ctx context.Context, limit int)) *Orderer_ListProcessing_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *Orderer_ListProcessing_Call) Return(_a0 []models.Order, _a1 error) *Orderer_ListProcessing_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Orderer_ListProcessing_Call) RunAndReturn(run func(context.Context, int) ([]models.Order, error)) *Orderer_ListProcessing_Call {
	_c.Call.Return(run)
	return _c
}

// OrderByID provides a mock function with given fields: _a0, _a1
func (_m *Orderer) OrderByID(_a0 context.Context, _a1 models.OrderID) (models.Order, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for OrderByID")
	}

	var r0 models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.OrderID) (models.Order, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.OrderID) models.Order); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(models.Order)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.OrderID) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Orderer_OrderByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OrderByID'
type Orderer_OrderByID_Call struct {
	*mock.Call
}

// OrderByID is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 models.OrderID
func (_e *Orderer_Expecter) OrderByID(_a0 interface{}, _a1 interface{}) *Orderer_OrderByID_Call {
	return &Orderer_OrderByID_Call{Call: _e.mock.On("OrderByID", _a0, _a1)}
}

func (_c *Orderer_OrderByID_Call) Run(run func(_a0 context.Context, _a1 models.OrderID)) *Orderer_OrderByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.OrderID))
	})
	return _c
}

func (_c *Orderer_OrderByID_Call) Return(_a0 models.Order, _a1 error) *Orderer_OrderByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Orderer_OrderByID_Call) RunAndReturn(run func(context.Context, models.OrderID) (models.Order, error)) *Orderer_OrderByID_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatus provides a mock function with given fields: _a0, _a1, _a2
func (_m *Orderer) UpdateStatus(_a0 context.Context, _a1 models.OrderID, _a2 models.OrderStatus) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.OrderID, models.OrderStatus) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Orderer_UpdateStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatus'
type Orderer_UpdateStatus_Call struct {
	*mock.Call
}

// UpdateStatus is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 models.OrderID
//   - _a2 models.OrderStatus
func (_e *Orderer_Expecter) UpdateStatus(_a0 interface{}, _a1 interface{}, _a2 interface{}) *Orderer_UpdateStatus_Call {
	return &Orderer_UpdateStatus_Call{Call: _e.mock.On("UpdateStatus", _a0, _a1, _a2)}
}

func (_c *Orderer_UpdateStatus_Call) Run(run func(_a0 context.Context, _a1 models.OrderID, _a2 models.OrderStatus)) *Orderer_UpdateStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.OrderID), args[2].(models.OrderStatus))
	})
	return _c
}

func (_c *Orderer_UpdateStatus_Call) Return(_a0 error) *Orderer_UpdateStatus_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Orderer_UpdateStatus_Call) RunAndReturn(run func(context.Context, models.OrderID, models.OrderStatus) error) *Orderer_UpdateStatus_Call {
	_c.Call.Return(run)
	return _c
}

// ValidateNumLuhn provides a mock function with given fields: num
func (_m *Orderer) ValidateNumLuhn(num string) error {
	ret := _m.Called(num)

	if len(ret) == 0 {
		panic("no return value specified for ValidateNumLuhn")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(num)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Orderer_ValidateNumLuhn_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateNumLuhn'
type Orderer_ValidateNumLuhn_Call struct {
	*mock.Call
}

// ValidateNumLuhn is a helper method to define mock.On call
//   - num string
func (_e *Orderer_Expecter) ValidateNumLuhn(num interface{}) *Orderer_ValidateNumLuhn_Call {
	return &Orderer_ValidateNumLuhn_Call{Call: _e.mock.On("ValidateNumLuhn", num)}
}

func (_c *Orderer_ValidateNumLuhn_Call) Run(run func(num string)) *Orderer_ValidateNumLuhn_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Orderer_ValidateNumLuhn_Call) Return(_a0 error) *Orderer_ValidateNumLuhn_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Orderer_ValidateNumLuhn_Call) RunAndReturn(run func(string) error) *Orderer_ValidateNumLuhn_Call {
	_c.Call.Return(run)
	return _c
}

// NewOrderer creates a new instance of Orderer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrderer(t interface {
	mock.TestingT
	Cleanup(func())
}) *Orderer {
	mock := &Orderer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
