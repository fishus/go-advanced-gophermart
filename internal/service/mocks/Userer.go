// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/fishus/go-advanced-gophermart/pkg/models"
	mock "github.com/stretchr/testify/mock"

	user "github.com/fishus/go-advanced-gophermart/internal/service/user"
)

// Userer is an autogenerated mock type for the Userer type
type Userer struct {
	mock.Mock
}

type Userer_Expecter struct {
	mock *mock.Mock
}

func (_m *Userer) EXPECT() *Userer_Expecter {
	return &Userer_Expecter{mock: &_m.Mock}
}

// BuildToken provides a mock function with given fields: _a0
func (_m *Userer) BuildToken(_a0 models.UserID) (string, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for BuildToken")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(models.UserID) (string, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(models.UserID) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(models.UserID) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Userer_BuildToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BuildToken'
type Userer_BuildToken_Call struct {
	*mock.Call
}

// BuildToken is a helper method to define mock.On call
//   - _a0 models.UserID
func (_e *Userer_Expecter) BuildToken(_a0 interface{}) *Userer_BuildToken_Call {
	return &Userer_BuildToken_Call{Call: _e.mock.On("BuildToken", _a0)}
}

func (_c *Userer_BuildToken_Call) Run(run func(_a0 models.UserID)) *Userer_BuildToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(models.UserID))
	})
	return _c
}

func (_c *Userer_BuildToken_Call) Return(_a0 string, _a1 error) *Userer_BuildToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Userer_BuildToken_Call) RunAndReturn(run func(models.UserID) (string, error)) *Userer_BuildToken_Call {
	_c.Call.Return(run)
	return _c
}

// CheckAuthorizationHeader provides a mock function with given fields: auth
func (_m *Userer) CheckAuthorizationHeader(auth string) (*user.JWTClaims, error) {
	ret := _m.Called(auth)

	if len(ret) == 0 {
		panic("no return value specified for CheckAuthorizationHeader")
	}

	var r0 *user.JWTClaims
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*user.JWTClaims, error)); ok {
		return rf(auth)
	}
	if rf, ok := ret.Get(0).(func(string) *user.JWTClaims); ok {
		r0 = rf(auth)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.JWTClaims)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(auth)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Userer_CheckAuthorizationHeader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckAuthorizationHeader'
type Userer_CheckAuthorizationHeader_Call struct {
	*mock.Call
}

// CheckAuthorizationHeader is a helper method to define mock.On call
//   - auth string
func (_e *Userer_Expecter) CheckAuthorizationHeader(auth interface{}) *Userer_CheckAuthorizationHeader_Call {
	return &Userer_CheckAuthorizationHeader_Call{Call: _e.mock.On("CheckAuthorizationHeader", auth)}
}

func (_c *Userer_CheckAuthorizationHeader_Call) Run(run func(auth string)) *Userer_CheckAuthorizationHeader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Userer_CheckAuthorizationHeader_Call) Return(_a0 *user.JWTClaims, _a1 error) *Userer_CheckAuthorizationHeader_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Userer_CheckAuthorizationHeader_Call) RunAndReturn(run func(string) (*user.JWTClaims, error)) *Userer_CheckAuthorizationHeader_Call {
	_c.Call.Return(run)
	return _c
}

// DecryptToken provides a mock function with given fields: tokenString
func (_m *Userer) DecryptToken(tokenString string) (*user.JWTClaims, error) {
	ret := _m.Called(tokenString)

	if len(ret) == 0 {
		panic("no return value specified for DecryptToken")
	}

	var r0 *user.JWTClaims
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*user.JWTClaims, error)); ok {
		return rf(tokenString)
	}
	if rf, ok := ret.Get(0).(func(string) *user.JWTClaims); ok {
		r0 = rf(tokenString)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.JWTClaims)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tokenString)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Userer_DecryptToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DecryptToken'
type Userer_DecryptToken_Call struct {
	*mock.Call
}

// DecryptToken is a helper method to define mock.On call
//   - tokenString string
func (_e *Userer_Expecter) DecryptToken(tokenString interface{}) *Userer_DecryptToken_Call {
	return &Userer_DecryptToken_Call{Call: _e.mock.On("DecryptToken", tokenString)}
}

func (_c *Userer_DecryptToken_Call) Run(run func(tokenString string)) *Userer_DecryptToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Userer_DecryptToken_Call) Return(_a0 *user.JWTClaims, _a1 error) *Userer_DecryptToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Userer_DecryptToken_Call) RunAndReturn(run func(string) (*user.JWTClaims, error)) *Userer_DecryptToken_Call {
	_c.Call.Return(run)
	return _c
}

// Login provides a mock function with given fields: _a0, _a1
func (_m *Userer) Login(_a0 context.Context, _a1 models.User) (models.UserID, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Login")
	}

	var r0 models.UserID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.User) (models.UserID, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.User) models.UserID); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(models.UserID)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.User) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Userer_Login_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Login'
type Userer_Login_Call struct {
	*mock.Call
}

// Login is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 models.User
func (_e *Userer_Expecter) Login(_a0 interface{}, _a1 interface{}) *Userer_Login_Call {
	return &Userer_Login_Call{Call: _e.mock.On("Login", _a0, _a1)}
}

func (_c *Userer_Login_Call) Run(run func(_a0 context.Context, _a1 models.User)) *Userer_Login_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.User))
	})
	return _c
}

func (_c *Userer_Login_Call) Return(_a0 models.UserID, _a1 error) *Userer_Login_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Userer_Login_Call) RunAndReturn(run func(context.Context, models.User) (models.UserID, error)) *Userer_Login_Call {
	_c.Call.Return(run)
	return _c
}

// LoyaltyAddWithdraw provides a mock function with given fields: ctx, userID, orderNum, withdraw
func (_m *Userer) LoyaltyAddWithdraw(ctx context.Context, userID models.UserID, orderNum string, withdraw float64) error {
	ret := _m.Called(ctx, userID, orderNum, withdraw)

	if len(ret) == 0 {
		panic("no return value specified for LoyaltyAddWithdraw")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID, string, float64) error); ok {
		r0 = rf(ctx, userID, orderNum, withdraw)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Userer_LoyaltyAddWithdraw_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LoyaltyAddWithdraw'
type Userer_LoyaltyAddWithdraw_Call struct {
	*mock.Call
}

// LoyaltyAddWithdraw is a helper method to define mock.On call
//   - ctx context.Context
//   - userID models.UserID
//   - orderNum string
//   - withdraw float64
func (_e *Userer_Expecter) LoyaltyAddWithdraw(ctx interface{}, userID interface{}, orderNum interface{}, withdraw interface{}) *Userer_LoyaltyAddWithdraw_Call {
	return &Userer_LoyaltyAddWithdraw_Call{Call: _e.mock.On("LoyaltyAddWithdraw", ctx, userID, orderNum, withdraw)}
}

func (_c *Userer_LoyaltyAddWithdraw_Call) Run(run func(ctx context.Context, userID models.UserID, orderNum string, withdraw float64)) *Userer_LoyaltyAddWithdraw_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.UserID), args[2].(string), args[3].(float64))
	})
	return _c
}

func (_c *Userer_LoyaltyAddWithdraw_Call) Return(_a0 error) *Userer_LoyaltyAddWithdraw_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Userer_LoyaltyAddWithdraw_Call) RunAndReturn(run func(context.Context, models.UserID, string, float64) error) *Userer_LoyaltyAddWithdraw_Call {
	_c.Call.Return(run)
	return _c
}

// LoyaltyUserBalance provides a mock function with given fields: _a0, _a1
func (_m *Userer) LoyaltyUserBalance(_a0 context.Context, _a1 models.UserID) (models.LoyaltyBalance, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for LoyaltyUserBalance")
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

// Userer_LoyaltyUserBalance_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LoyaltyUserBalance'
type Userer_LoyaltyUserBalance_Call struct {
	*mock.Call
}

// LoyaltyUserBalance is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 models.UserID
func (_e *Userer_Expecter) LoyaltyUserBalance(_a0 interface{}, _a1 interface{}) *Userer_LoyaltyUserBalance_Call {
	return &Userer_LoyaltyUserBalance_Call{Call: _e.mock.On("LoyaltyUserBalance", _a0, _a1)}
}

func (_c *Userer_LoyaltyUserBalance_Call) Run(run func(_a0 context.Context, _a1 models.UserID)) *Userer_LoyaltyUserBalance_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.UserID))
	})
	return _c
}

func (_c *Userer_LoyaltyUserBalance_Call) Return(_a0 models.LoyaltyBalance, _a1 error) *Userer_LoyaltyUserBalance_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Userer_LoyaltyUserBalance_Call) RunAndReturn(run func(context.Context, models.UserID) (models.LoyaltyBalance, error)) *Userer_LoyaltyUserBalance_Call {
	_c.Call.Return(run)
	return _c
}

// LoyaltyUserWithdrawals provides a mock function with given fields: _a0, _a1
func (_m *Userer) LoyaltyUserWithdrawals(_a0 context.Context, _a1 models.UserID) ([]models.LoyaltyHistory, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for LoyaltyUserWithdrawals")
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

// Userer_LoyaltyUserWithdrawals_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LoyaltyUserWithdrawals'
type Userer_LoyaltyUserWithdrawals_Call struct {
	*mock.Call
}

// LoyaltyUserWithdrawals is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 models.UserID
func (_e *Userer_Expecter) LoyaltyUserWithdrawals(_a0 interface{}, _a1 interface{}) *Userer_LoyaltyUserWithdrawals_Call {
	return &Userer_LoyaltyUserWithdrawals_Call{Call: _e.mock.On("LoyaltyUserWithdrawals", _a0, _a1)}
}

func (_c *Userer_LoyaltyUserWithdrawals_Call) Run(run func(_a0 context.Context, _a1 models.UserID)) *Userer_LoyaltyUserWithdrawals_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.UserID))
	})
	return _c
}

func (_c *Userer_LoyaltyUserWithdrawals_Call) Return(_a0 []models.LoyaltyHistory, _a1 error) *Userer_LoyaltyUserWithdrawals_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Userer_LoyaltyUserWithdrawals_Call) RunAndReturn(run func(context.Context, models.UserID) ([]models.LoyaltyHistory, error)) *Userer_LoyaltyUserWithdrawals_Call {
	_c.Call.Return(run)
	return _c
}

// Register provides a mock function with given fields: _a0, _a1
func (_m *Userer) Register(_a0 context.Context, _a1 models.User) (models.UserID, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Register")
	}

	var r0 models.UserID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.User) (models.UserID, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.User) models.UserID); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(models.UserID)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.User) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Userer_Register_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Register'
type Userer_Register_Call struct {
	*mock.Call
}

// Register is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 models.User
func (_e *Userer_Expecter) Register(_a0 interface{}, _a1 interface{}) *Userer_Register_Call {
	return &Userer_Register_Call{Call: _e.mock.On("Register", _a0, _a1)}
}

func (_c *Userer_Register_Call) Run(run func(_a0 context.Context, _a1 models.User)) *Userer_Register_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.User))
	})
	return _c
}

func (_c *Userer_Register_Call) Return(_a0 models.UserID, _a1 error) *Userer_Register_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Userer_Register_Call) RunAndReturn(run func(context.Context, models.User) (models.UserID, error)) *Userer_Register_Call {
	_c.Call.Return(run)
	return _c
}

// UserByID provides a mock function with given fields: _a0, _a1
func (_m *Userer) UserByID(_a0 context.Context, _a1 models.UserID) (models.User, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for UserByID")
	}

	var r0 models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID) (models.User, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID) models.User); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.UserID) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Userer_UserByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UserByID'
type Userer_UserByID_Call struct {
	*mock.Call
}

// UserByID is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 models.UserID
func (_e *Userer_Expecter) UserByID(_a0 interface{}, _a1 interface{}) *Userer_UserByID_Call {
	return &Userer_UserByID_Call{Call: _e.mock.On("UserByID", _a0, _a1)}
}

func (_c *Userer_UserByID_Call) Run(run func(_a0 context.Context, _a1 models.UserID)) *Userer_UserByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.UserID))
	})
	return _c
}

func (_c *Userer_UserByID_Call) Return(_a0 models.User, _a1 error) *Userer_UserByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Userer_UserByID_Call) RunAndReturn(run func(context.Context, models.UserID) (models.User, error)) *Userer_UserByID_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserer creates a new instance of Userer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserer(t interface {
	mock.TestingT
	Cleanup(func())
}) *Userer {
	mock := &Userer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
