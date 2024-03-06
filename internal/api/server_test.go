package api

import (
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"

	"github.com/fishus/go-advanced-gophermart/internal/api/mocks"
	"github.com/fishus/go-advanced-gophermart/internal/service"
)

type APITestSuite struct {
	suite.Suite
	server  *server
	srv     *httptest.Server
	client  *resty.Client
	service *mocks.Servicer
	loyalty *mocks.AccrualDaemon
}

func (ts *APITestSuite) SetupSuite() {
	ts.service = mocks.NewServicer(ts.T())
	ts.loyalty = mocks.NewAccrualDaemon(ts.T())

	ts.server = &server{
		cfg:     &Config{},
		service: ts.service,
		loyalty: ts.loyalty,
	}

	ts.srv = httptest.NewServer(Router(ts.server))
	ts.server.cfg.ServerAddr = ts.srv.URL
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

func (ts *APITestSuite) setService(o service.Orderer, u service.Userer) {
	if o != nil {
		ts.service.EXPECT().Order().Return(o)
	}
	if u != nil {
		ts.service.EXPECT().User().Return(u)
	}
}

func TestApi(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}