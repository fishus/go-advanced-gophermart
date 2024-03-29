// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/fishus/go-advanced-gophermart/pkg/models"
	mock "github.com/stretchr/testify/mock"
)

// AccrualDaemon is an autogenerated mock type for the AccrualDaemon type
type AccrualDaemon struct {
	mock.Mock
}

type AccrualDaemon_Expecter struct {
	mock *mock.Mock
}

func (_m *AccrualDaemon) EXPECT() *AccrualDaemon_Expecter {
	return &AccrualDaemon_Expecter{mock: &_m.Mock}
}

// AddNewOrder provides a mock function with given fields: _a0, _a1
func (_m *AccrualDaemon) AddNewOrder(_a0 context.Context, _a1 models.Order) {
	_m.Called(_a0, _a1)
}

// AccrualDaemon_AddNewOrder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddNewOrder'
type AccrualDaemon_AddNewOrder_Call struct {
	*mock.Call
}

// AddNewOrder is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 models.Order
func (_e *AccrualDaemon_Expecter) AddNewOrder(_a0 interface{}, _a1 interface{}) *AccrualDaemon_AddNewOrder_Call {
	return &AccrualDaemon_AddNewOrder_Call{Call: _e.mock.On("AddNewOrder", _a0, _a1)}
}

func (_c *AccrualDaemon_AddNewOrder_Call) Run(run func(_a0 context.Context, _a1 models.Order)) *AccrualDaemon_AddNewOrder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.Order))
	})
	return _c
}

func (_c *AccrualDaemon_AddNewOrder_Call) Return() *AccrualDaemon_AddNewOrder_Call {
	_c.Call.Return()
	return _c
}

func (_c *AccrualDaemon_AddNewOrder_Call) RunAndReturn(run func(context.Context, models.Order)) *AccrualDaemon_AddNewOrder_Call {
	_c.Call.Return(run)
	return _c
}

// NewAccrualDaemon creates a new instance of AccrualDaemon. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAccrualDaemon(t interface {
	mock.TestingT
	Cleanup(func())
}) *AccrualDaemon {
	mock := &AccrualDaemon{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
