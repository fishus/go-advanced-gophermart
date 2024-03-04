package accrual

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

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
					if err != nil {
						switch {
						case errors.Is(err, ErrIsShutdown),
							errors.Is(err, ErrMaxAttemptsReached),
							errors.Is(err, ErrOrderNotRegistered),
							errors.Is(err, ErrAPIServerError):
							err = d.service.Order().UpdateStatus(ctx, order.ID, models.OrderStatusProcessing)
							if err != nil {
								logger.Log.Error(err.Error(), logger.String("OrderID", order.ID.String()))
							}

						default:
							err = d.service.Order().UpdateStatus(ctx, order.ID, models.OrderStatusInvalid)
							if err != nil {
								logger.Log.Error(err.Error(), logger.String("OrderID", order.ID.String()))
							}
						}
						return
					}

					if acc.Num != order.Num {
						err = d.service.Order().UpdateStatus(ctx, order.ID, models.OrderStatusInvalid)
						if err != nil {
							logger.Log.Error(err.Error(), logger.String("OrderID", order.ID.String()))
						}
						return
					}

					switch acc.Status {
					case OrderAccrualStatusInvalid:
						err = d.service.Order().UpdateStatus(ctx, order.ID, models.OrderStatusInvalid)
						if err != nil {
							logger.Log.Error(err.Error(), logger.String("OrderID", order.ID.String()))
						}
					case OrderAccrualStatusRegistered,
						OrderAccrualStatusProcessing:
						err = d.service.Order().UpdateStatus(ctx, order.ID, models.OrderStatusProcessing)
						if err != nil {
							logger.Log.Error(err.Error(), logger.String("OrderID", order.ID.String()))
						}
					case OrderAccrualStatusProcessed:
						err = d.service.Order().AddAccrual(ctx, order.ID, acc.Accrual)
						if err != nil {
							logger.Log.Error(err.Error(), logger.String("OrderID", order.ID.String()))
						}
					}
				}()
			}
		}
	}()
}

// requestOrderAccrual makes a request for the status of accrual loyalty points for an order
func (d *daemon) requestOrderAccrual(ctx context.Context, num string) (*OrderAccrual, error) {
	if num == "" {
		return nil, errors.New("order number is empty")
	}

	i := 0 // Attempts counter
	var respOrder OrderAccrual
	for {
		select {
		case <-ctx.Done():
			return nil, ErrIsShutdown
		case <-d.chShutdown:
			return nil, ErrIsShutdown
		default:
		}

		i++
		if i != 0 && i > d.cfg.MaxAttempts {
			return nil, ErrMaxAttemptsReached
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
			return nil, err
		}

		switch resp.StatusCode() {
		// успешная обработка запроса
		case http.StatusOK:
			if err = respOrder.Status.Validate(); err != nil {
				return nil, err
			}
			return &respOrder, nil

		// превышено количество запросов к сервису
		case http.StatusTooManyRequests:
			retryAfter, err := strconv.Atoi(resp.Header().Get("Retry-After"))
			if err != nil {
				return nil, err
			}
			d.doDelay(ctx, time.Duration(retryAfter+1)*time.Second)
			continue

		// заказ не зарегистрирован в системе расчёта
		case http.StatusNoContent:
			return nil, ErrOrderNotRegistered

		// внутренняя ошибка сервера
		case http.StatusInternalServerError:
			return nil, ErrAPIServerError

		default:
			return nil, errors.New(http.StatusText(resp.StatusCode()))
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
