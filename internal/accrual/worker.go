package accrual

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/app/config"
	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

var (
	ErrIsShutdown         = errors.New("service is shutting down")
	ErrEmptyOrderNum      = errors.New("order number is empty")
	ErrAPIServerError     = errors.New("accrual api server error")
	ErrOrderNotRegistered = errors.New("the order is not registered in the system")
	ErrMaxAttemptsReached = errors.New("maximum number of attempts has been reached")
)

type orderAccrualError struct {
	err         error
	orderStatus models.OrderStatus
}

func (e *orderAccrualError) Error() string {
	return fmt.Sprintf("%v %v", e.orderStatus, e.err)
}

func (e *orderAccrualError) Unwrap() error {
	return e.err
}

func newOrderAccrualError(err error, s models.OrderStatus) error {
	return &orderAccrualError{
		err:         err,
		orderStatus: s,
	}
}

// workerRequestOrderAccrual the worker sends requests for information about accrual loyalty points and processes the responses
func (d *daemon) workerRequestOrderAccrual(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				select {
				case <-d.chOrders: // read from chan before exit
				default:
				}
				return
			case <-d.chShutdown:
				select {
				case <-d.chOrders: // read from chan before exit
				default:
				}
				return

			case order, opened := <-d.chOrders:
				if !opened || ctx.Err() != nil {
					return
				}
				d.wg.Add(1)
				func() {
					defer d.wg.Done()
					acc, err := d.requestOrderAccrual(ctx, order.Num)

					if err == nil && acc.Num != order.Num {
						err = newOrderAccrualError(errors.New("order number mismatch between request and response"), models.OrderStatusInvalid)
					}

					if err != nil {
						var e *orderAccrualError
						if errors.As(err, &e) {
							err = d.service.Order().UpdateStatus(ctx, order.ID, e.orderStatus)
							if err != nil {
								logger.Log.Error(err.Error(), logger.String("OrderID", order.ID.String()))
							}
						}
						return
					}

					switch acc.Status {
					case OrderAccrualStatusInvalid:
						err = d.service.Order().UpdateStatus(ctx, order.ID, models.OrderStatusInvalid)
					case OrderAccrualStatusRegistered,
						OrderAccrualStatusProcessing:
						err = d.service.Order().UpdateStatus(ctx, order.ID, models.OrderStatusProcessing)
					case OrderAccrualStatusProcessed:
						err = d.service.Order().AddAccrual(ctx, order.ID, acc.Accrual)
					}

					if err != nil {
						logger.Log.Error(err.Error(), logger.String("OrderID", order.ID.String()))
					}
				}()
			}
		}
	}()
}

// requestOrderAccrual makes a request for the status of accrual loyalty points for an order
func (d *daemon) requestOrderAccrual(ctx context.Context, num string) (*OrderAccrual, error) {
	if num == "" {
		return nil, newOrderAccrualError(ErrEmptyOrderNum, models.OrderStatusInvalid)
	}

	i := 0 // Attempts counter
	var respOrder OrderAccrual
	for {
		select {
		case <-ctx.Done():
			return nil, newOrderAccrualError(ErrIsShutdown, models.OrderStatusProcessing)
		case <-d.chShutdown:
			return nil, newOrderAccrualError(ErrIsShutdown, models.OrderStatusProcessing)
		default:
		}

		i++
		if i != 0 && i > d.cfg.MaxAttempts {
			return nil, newOrderAccrualError(ErrMaxAttemptsReached, models.OrderStatusProcessing)
		}

		d.delayCond.L.Lock()
		for d.delay.Load() {
			d.delayCond.Wait()
		}
		d.delayCond.L.Unlock()

		req := d.client.R().
			SetContext(ctx).
			SetResult(&respOrder).
			SetPathParams(map[string]string{
				"number": num,
			})
		resp, err := req.Get("/api/orders/{number}")
		if err != nil {
			logger.Log.Error(err.Error(), logger.String("OrderNum", num))
			return nil, newOrderAccrualError(err, models.OrderStatusProcessing)
		}

		switch resp.StatusCode() {
		// успешная обработка запроса
		case http.StatusOK:
			if err = respOrder.Status.Validate(); err != nil {
				return nil, newOrderAccrualError(err, models.OrderStatusInvalid)
			}
			respOrder.Accrual = respOrder.Accrual.Round(config.DecimalExponent)
			return &respOrder, nil

		// превышено количество запросов к сервису
		case http.StatusTooManyRequests:
			retryAfter, err := strconv.Atoi(resp.Header().Get("Retry-After"))
			if err != nil {
				return nil, newOrderAccrualError(err, models.OrderStatusInvalid)
			}
			d.doDelay(ctx, time.Duration(retryAfter+1)*time.Second)
			continue

		// заказ не зарегистрирован в системе расчёта
		case http.StatusNoContent:
			return nil, newOrderAccrualError(ErrOrderNotRegistered, models.OrderStatusProcessing)

		// внутренняя ошибка сервера
		case http.StatusInternalServerError:
			return nil, newOrderAccrualError(ErrAPIServerError, models.OrderStatusProcessing)

		default:
			return nil, newOrderAccrualError(errors.New(http.StatusText(resp.StatusCode())), models.OrderStatusInvalid)
		}
	}
}

func (d *daemon) doDelay(ctx context.Context, delayDuration time.Duration) {
	d.delayCond.L.Lock()
	d.delay.Store(true)

	expire := time.Now().Add(delayDuration)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
exit:
	for {
		select {
		case <-ctx.Done():
			// interrupt the pause early because the application is shutting down
			break exit
		case <-d.chShutdown:
			// interrupt the pause early because the application is shutting down
			break exit
		default:
		}

		t := <-ticker.C
		if t.Equal(expire) || t.After(expire) {
			break
		}
	}

	d.delay.Store(false)
	d.delayCond.Broadcast()
	d.delayCond.L.Unlock()
}
