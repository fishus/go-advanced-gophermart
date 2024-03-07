package order

import (
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"

	"github.com/fishus/go-advanced-gophermart/internal/api/mocks"
	"github.com/fishus/go-advanced-gophermart/internal/service"
)

type APITestSuite struct {
	suite.Suite
	api     *api
	srv     *httptest.Server
	client  *resty.Client
	service *mocks.Servicer
	daemon  *mocks.AccrualDaemon
}

func (ts *APITestSuite) SetupSuite() {
	ts.service = mocks.NewServicer(ts.T())
	ts.daemon = mocks.NewAccrualDaemon(ts.T())

	ts.api = &api{
		service: ts.service,
		daemon:  ts.daemon,
	}

	ts.srv = httptest.NewServer(ts.api.router())
	ts.client = resty.New().SetBaseURL(ts.srv.URL)
}

func (ts *APITestSuite) TearDownSuite() {
	ts.srv.Close()
}

func (ts *APITestSuite) TearDownTest() {
	// Reset mock
	for _, mc := range ts.service.ExpectedCalls {
		mc.Unset()
	}
}

func (ts *APITestSuite) TearDownSubTest() {
	// Reset mock
	for _, mc := range ts.service.ExpectedCalls {
		mc.Unset()
	}
}

func (ts *APITestSuite) setServiceOrder(s service.Orderer) {
	ts.service.EXPECT().Order().Return(s)
}

func (ts *APITestSuite) setServiceUser(s service.Userer) {
	ts.service.EXPECT().User().Return(s)
}

func TestApi(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (a *api) router() chi.Router {
	r := chi.NewRouter()

	r.Post("/add", a.Add)  // Загрузка номера заказа для расчёта
	r.Get("/list", a.List) // Список загруженных номеров заказов

	return r
}
