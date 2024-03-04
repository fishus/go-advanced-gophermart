package accrual

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"
)

func (ts *LoyaltyTestSuite) TestRequestOrderAccrual() {
	ctx := context.Background()

	type orderResp struct {
		Num     string             `json:"order"`
		Status  OrderAccrualStatus `json:"status"`
		Accrual float64            `json:"accrual"`
	}

	// заказ зарегистрирован, но вознаграждение не рассчитано
	ts.Run("Status 200: Order REGISTERED", func() {
		orderNum := "4997753043"

		wantOrder := &orderResp{
			Num:     orderNum,
			Status:  "REGISTERED",
			Accrual: 0,
		}

		ts.fakeLoyaltyAPI.contentType = "application/json; charset=utf-8"
		ts.fakeLoyaltyAPI.statusCode = http.StatusOK

		resp, err := json.Marshal(wantOrder)
		ts.Require().NoError(err)
		ts.fakeLoyaltyAPI.resp = string(resp)

		order, err := ts.daemon.requestOrderAccrual(ctx, orderNum)
		ts.NoError(err)
		ts.EqualValues(wantOrder, order)
		ts.Equal(orderNum, ts.fakeLoyaltyAPI.actual.number)
	})

	// заказ не принят к расчёту, и вознаграждение не будет начислено
	ts.Run("Status 200: Order INVALID", func() {
		orderNum := "9400781309"

		wantOrder := &orderResp{
			Num:     orderNum,
			Status:  "INVALID",
			Accrual: 0,
		}

		ts.fakeLoyaltyAPI.contentType = "application/json; charset=utf-8"
		ts.fakeLoyaltyAPI.statusCode = http.StatusOK

		resp, err := json.Marshal(wantOrder)
		ts.Require().NoError(err)
		ts.fakeLoyaltyAPI.resp = string(resp)

		order, err := ts.daemon.requestOrderAccrual(ctx, orderNum)
		ts.NoError(err)
		ts.EqualValues(wantOrder, order)
		ts.Equal(orderNum, ts.fakeLoyaltyAPI.actual.number)
	})

	// расчёт начисления в процессе;
	ts.Run("Status 200: Order PROCESSING", func() {
		orderNum := "8163091187"

		wantOrder := &orderResp{
			Num:     orderNum,
			Status:  "PROCESSING",
			Accrual: 0,
		}

		ts.fakeLoyaltyAPI.contentType = "application/json; charset=utf-8"
		ts.fakeLoyaltyAPI.statusCode = http.StatusOK

		resp, err := json.Marshal(wantOrder)
		ts.Require().NoError(err)
		ts.fakeLoyaltyAPI.resp = string(resp)

		order, err := ts.daemon.requestOrderAccrual(ctx, orderNum)
		ts.NoError(err)
		ts.EqualValues(wantOrder, order)
		ts.Equal(orderNum, ts.fakeLoyaltyAPI.actual.number)
	})

	// расчёт начисления окончен
	ts.Run("Status 200: Order PROCESSED", func() {
		orderNum := "5347676263"

		wantOrder := &orderResp{
			Num:     orderNum,
			Status:  "PROCESSED",
			Accrual: 123.456,
		}

		ts.fakeLoyaltyAPI.contentType = "application/json; charset=utf-8"
		ts.fakeLoyaltyAPI.statusCode = http.StatusOK

		resp, err := json.Marshal(wantOrder)
		ts.Require().NoError(err)
		ts.fakeLoyaltyAPI.resp = string(resp)

		order, err := ts.daemon.requestOrderAccrual(ctx, orderNum)
		ts.NoError(err)
		ts.EqualValues(wantOrder, order)
		ts.Equal(orderNum, ts.fakeLoyaltyAPI.actual.number)
	})

	// заказ не зарегистрирован в системе расчёта
	ts.Run("Status 204", func() {
		orderNum := "3903733214"

		ts.fakeLoyaltyAPI.statusCode = http.StatusNoContent

		_, err := ts.daemon.requestOrderAccrual(ctx, orderNum)
		ts.Error(err)
		ts.ErrorIs(err, ErrOrderNotRegistered)
		ts.Equal(orderNum, ts.fakeLoyaltyAPI.actual.number)
	})

	// внутренняя ошибка сервера
	ts.Run("Status 500", func() {
		orderNum := "0306558669"

		ts.fakeLoyaltyAPI.statusCode = http.StatusInternalServerError

		_, err := ts.daemon.requestOrderAccrual(ctx, orderNum)
		ts.Error(err)
		ts.ErrorIs(err, ErrAPIServerError)
		ts.Equal(orderNum, ts.fakeLoyaltyAPI.actual.number)

		ts.fakeLoyaltyAPI.retryAfter = 1 // Delay is 1 sec
	})

	// превышено количество запросов к сервису
	ts.Run("Status 429", func() {
		orderNum := "1853241857"

		// After several attempts we return status 204
		ts.fakeLoyaltyAPI.statusCode = http.StatusNoContent
		ts.fakeLoyaltyAPI.retryAfter = 1 // Delay is 1 sec
		ts.fakeLoyaltyAPI.retries = 2

		start := time.Now()
		_, err := ts.daemon.requestOrderAccrual(ctx, orderNum)
		ts.Error(err)
		ts.ErrorIs(err, ErrOrderNotRegistered)
		ts.Equal(orderNum, ts.fakeLoyaltyAPI.actual.number)
		ts.Equal(ts.fakeLoyaltyAPI.retries, ts.fakeLoyaltyAPI.actual.retries)
		ts.Greater(time.Since(start), (time.Duration(ts.fakeLoyaltyAPI.retryAfter*ts.fakeLoyaltyAPI.retries) * time.Second))
	})
}

// Test handler: Тестовый расчёт начислений баллов лояльности
func (s *fakeLoyaltyAPI) orderAccrual(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	statusCode := s.statusCode
	contentType := s.contentType
	resp := s.resp

	s.actual.number = chi.URLParam(r, "number")

	if s.retries > 0 {
		s.actual.retries++
	}

	if s.retryAfter > 0 && s.retries > 0 && s.actual.retries < s.retries {
		statusCode = http.StatusTooManyRequests
		contentType = "text/plain"
		i := strconv.Itoa(s.retryAfter)
		w.Header().Set("Retry-After", i)
		resp = "No more than 60 requests per minute allowed"
	}

	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.WriteHeader(statusCode)

	if resp != "" {
		_, _ = io.WriteString(w, resp)
	}
}

func (ts *LoyaltyTestSuite) TestDoDelay() {
	timeout := 2 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var execTime atomic.Int64

	start := time.Now()

	go func() {
		ts.daemon.doDelay(ctx, 1*time.Second)
	}()

	go func() {
		ts.daemon.delayCond.L.Lock()
		ts.daemon.delayCond.Wait()
		ts.daemon.delayCond.L.Unlock()
		if execTime.Load() == 0 {
			execTime.Store(int64(time.Since(start)))
			cancel()
		}
	}()

	<-ctx.Done()
	if execTime.Load() == 0 {
		execTime.Store(int64(time.Since(start)))
		ts.daemon.delayCond.Broadcast()
	}

	t := time.Duration(execTime.Load())
	ts.Less(t, timeout)
}
