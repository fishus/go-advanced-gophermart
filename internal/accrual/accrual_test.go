package accrual

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/suite"

	"github.com/fishus/go-advanced-gophermart/internal/accrual/mocks"
	"github.com/fishus/go-advanced-gophermart/internal/service"
)

type LoyaltyTestSuite struct {
	suite.Suite
	service       *mocks.Servicer
	daemon        *daemon
	loyaltyServer *httptest.Server
}

func (ts *LoyaltyTestSuite) SetupSuite() {
	ts.loyaltyServer = httptest.NewServer(ts.fakeLoyaltyRouter())

	cfg := &Config{
		APIAddr:        ts.loyaltyServer.URL,
		RequestTimeout: 5 * time.Second,
		MaxAttempts:    3,
		WorkersNum:     1,
	}

	ts.service = mocks.NewServicer(ts.T())
	ts.daemon = NewAccrual(cfg, ts.service)
}

func (ts *LoyaltyTestSuite) TearDownTest() {
	// Reset mock
	for _, mc := range ts.service.ExpectedCalls {
		mc.Unset()
	}
}

func (ts *LoyaltyTestSuite) TearDownSubTest() {
	// Reset mock
	for _, mc := range ts.service.ExpectedCalls {
		mc.Unset()
	}
}

func TestLoyalty(t *testing.T) {
	suite.Run(t, new(LoyaltyTestSuite))
}

func (ts *LoyaltyTestSuite) setService(o service.Orderer, u service.Userer) {
	if o != nil {
		ts.service.EXPECT().Order().Return(o)
	}
	if u != nil {
		ts.service.EXPECT().User().Return(u)
	}
}

func (ts *LoyaltyTestSuite) fakeLoyaltyRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	r.Get("/api/orders/{number}", ts.fakeLoyaltyOrderAccrual) // Тестовый расчёт начислений баллов лояльности

	return r
}

// Тестовый расчёт начислений баллов лояльности
func (ts *LoyaltyTestSuite) fakeLoyaltyOrderAccrual(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	orderNum := chi.URLParam(r, "number")

	// Читать параметры из заголовков

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 200 — успешная обработка запроса.
	/*
		{
		  "order": "<number>",
		  "status": "PROCESSED",
		  "accrual": 500
		}
		Поля объекта ответа:
		order — номер заказа;
		status — статус расчёта начисления:
		REGISTERED — заказ зарегистрирован, но вознаграждение не рассчитано;
		INVALID — заказ не принят к расчёту, и вознаграждение не будет начислено;
		PROCESSING — расчёт начисления в процессе;
		PROCESSED — расчёт начисления окончен;
		accrual — рассчитанные баллы к начислению, при отсутствии начисления — поле отсутствует в ответе.
	*/

	// 204 — заказ не зарегистрирован в системе расчёта.

	// 429 — превышено количество запросов к сервису.
	/*
	  429 Too Many Requests HTTP/1.1
	  Content-Type: text/plain
	  Retry-After: 60

	  No more than N requests per minute allowed
	*/

	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, orderNum)
}
