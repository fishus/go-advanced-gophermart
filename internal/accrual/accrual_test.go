package accrual

import (
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
	service        *mocks.Servicer
	daemon         *daemon
	loyaltyServer  *httptest.Server
	fakeLoyaltyAPI *fakeLoyaltyAPI
}

func (ts *LoyaltyTestSuite) SetupSuite() {
	ts.fakeLoyaltyAPI = newFakeLoyaltyAPI()
	ts.loyaltyServer = httptest.NewServer(ts.fakeLoyaltyAPI.Router())

	cfg := &Config{
		APIAddr:        ts.loyaltyServer.URL,
		RequestTimeout: 5 * time.Second,
		MaxAttempts:    3,
		WorkersNum:     1,
	}

	ts.service = mocks.NewServicer(ts.T())
	ts.daemon = NewAccrual(cfg, ts.service)
}

func (ts *LoyaltyTestSuite) TearDownSuite() {
	ts.loyaltyServer.Close()
	ts.daemon.Close()
}

func (ts *LoyaltyTestSuite) TearDownTest() {
	// Reset mock
	for _, mc := range ts.service.ExpectedCalls {
		mc.Unset()
	}
	ts.fakeLoyaltyAPI.Clear()
}

func (ts *LoyaltyTestSuite) TearDownSubTest() {
	// Reset mock
	for _, mc := range ts.service.ExpectedCalls {
		mc.Unset()
	}
	ts.fakeLoyaltyAPI.Clear()
}

func (ts *LoyaltyTestSuite) setService(o service.Orderer, u service.Userer) {
	if o != nil {
		ts.service.EXPECT().Order().Return(o)
	}
	if u != nil {
		ts.service.EXPECT().User().Return(u)
	}
}

func TestLoyalty(t *testing.T) {
	suite.Run(t, new(LoyaltyTestSuite))
}

type fakeLoyaltyAPI struct {
	retryAfter  int
	retries     int
	contentType string
	statusCode  int
	resp        string
	actual      struct {
		number  string
		retries int
	}
}

func (s *fakeLoyaltyAPI) Router() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Get("/api/orders/{number}", s.orderAccrual) // Тестовый расчёт начислений баллов лояльности
	return r
}

func (s *fakeLoyaltyAPI) Clear() {
	s.retryAfter = 0
	s.retries = 0
	s.contentType = "text/plain"
	s.statusCode = http.StatusOK
	s.resp = ""
	s.actual.number = ""
	s.actual.retries = 0
}

func newFakeLoyaltyAPI() *fakeLoyaltyAPI {
	api := &fakeLoyaltyAPI{}
	api.Clear()
	return api
}
