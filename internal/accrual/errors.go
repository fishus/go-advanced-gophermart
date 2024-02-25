package accrual

import (
	"errors"
)

var ErrAPIServerError = errors.New("accrual api server error")
var ErrOrderNotRegistered = errors.New("the order is not registered in the system")
var ErrIsShutdown = errors.New("service is shutting down")
var ErrMaxAttemptsReached = errors.New("maximum number of attempts has been reached")
