package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

// workerGetOrderAccrual the worker sends requests for information about accrual loyalty points and processes the responses
func (d *daemon) workerGetOrderAccrual(ctx context.Context, chOrders chan models.Order, chDelayed chan delayedOrder) {
	go func() {
		select {
		case <-ctx.Done():
			return
		default:
		}

		for order := range chOrders {
			ac, err := d.getOrderAccrual(ctx, order.Num)
			if err != nil {
				logger.Log.Error(err.Error(), logger.String("OrderNum", order.Num))

				var de *DelayedOrderError
				switch {
				case (errors.Is(err, ErrOrderNotRegistered) && errors.As(err, &de)),
					(errors.Is(err, ErrTooManyRequests) && errors.As(err, &de)),
					(errors.Is(err, ErrAPIServerError) && errors.As(err, &de)):
					chDelayed <- delayedOrder{order, time.Now().Add(de.Delay)}
					continue
				}

				// TODO перевести в статус INVALID
				continue
			}
			logger.Log.Info("OrderAccrual", logger.Any("OrderAccrual", ac))
			// TODO обработать статусы
			time.Sleep(100 * time.Millisecond)
		}

		// Статусы INVALID и PROCESSED являются окончательными.
	}()
}

// workerDelayed implements a queue of delayed orders
func (d *daemon) workerDelayed(ctx context.Context, chOrders chan<- models.Order) chan delayedOrder {
	chDelayed := make(chan delayedOrder, 10000)

	go func() {
		select {
		case <-ctx.Done():
			close(chDelayed)
			return
		default:
		}

		for delayedOrder := range chDelayed {
			if time.Now().After(delayedOrder.delay) {
				chOrders <- delayedOrder.order
			} else {
				chDelayed <- delayedOrder
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return chDelayed
}

// getOrderAccrual implements a request for the status of accrual loyalty points for an order
func (d *daemon) getOrderAccrual(ctx context.Context, num string) (*OrderAccrual, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if num == "" {
		return nil, errors.New("order number is empty")
	}

	req := d.client.R().
		SetContext(ctx).
		SetDoNotParseResponse(true).
		SetPathParams(map[string]string{
			"number": num,
		})
	resp, err := req.Get("/api/orders/{number}")
	if err != nil {
		return nil, err
	}
	rawBody := resp.RawBody()
	defer rawBody.Close()

	switch resp.StatusCode() {
	// успешная обработка запроса
	case http.StatusOK:
		// TODO

		var respData struct {
			Order   string
			Status  string
			Accrual float64 // FIXME
		}

		if err = json.NewDecoder(rawBody).Decode(&respData); err != nil {
			logger.Log.Debug(err.Error())
			return nil, err
		}

		status := OrderAccrualStatus(respData.Status)
		if err = status.Validate(); err != nil {
			return nil, err
		}

		ac := &OrderAccrual{
			Num:     respData.Order,
			Status:  status,
			Accrual: respData.Accrual,
		}

		return ac, nil

	// превышено количество запросов к сервису
	case http.StatusTooManyRequests:
		retryAfter, err := strconv.Atoi(resp.Header().Get("Retry-After"))
		if err != nil {
			return nil, err
		}
		return nil, NewDelayedOrderError(num, resp.StatusCode(), time.Duration(retryAfter)*time.Second, ErrTooManyRequests)
	// заказ не зарегистрирован в системе расчёта
	case http.StatusNoContent:
		return nil, NewDelayedOrderError(num, resp.StatusCode(), 2*time.Second, ErrOrderNotRegistered)
	// внутренняя ошибка сервера
	case http.StatusInternalServerError:
		return nil, NewDelayedOrderError(num, resp.StatusCode(), 2*time.Second, ErrAPIServerError)
	}

	return nil, NewDelayedOrderError(num, resp.StatusCode(), 0, errors.New(http.StatusText(resp.StatusCode())))
}
