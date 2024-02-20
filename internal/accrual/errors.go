package accrual

import (
	"errors"
	"fmt"
	"time"
)

var ErrAPIServerError = errors.New("accrual api server error")
var ErrOrderNotRegistered = errors.New("the order is not registered in the system")
var ErrTooManyRequests = errors.New("too many requests")

type DelayedOrderError struct {
	OrderNum   string
	StatusCode int
	Delay      time.Duration // Retry-After
	Err        error
}

func NewDelayedOrderError(num string, code int, delay time.Duration, err error) *DelayedOrderError {
	return &DelayedOrderError{
		OrderNum:   num,
		StatusCode: code,
		Delay:      delay,
		Err:        err,
	}
}

func (de *DelayedOrderError) Error() string {
	return fmt.Sprintf("[%d] #%s: %v", de.StatusCode, de.OrderNum, de.Err)
}

func (de *DelayedOrderError) Unwrap() error {
	return de.Err
}
