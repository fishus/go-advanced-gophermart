package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	oService "github.com/fishus/go-advanced-gophermart/internal/service/order"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	storageMocks "github.com/fishus/go-advanced-gophermart/internal/storage/mocks"
)

type UserServiceTestSuite struct {
	suite.Suite
	cfg     *Config
	storage *storageMocks.Storager
	*service
}

func (ts *UserServiceTestSuite) SetupSuite() {
	ts.cfg = &Config{
		JWTExpires:   15 * time.Minute,
		JWTSecretKey: "TestSecretKey",
	}
	ts.storage = storageMocks.NewStorager(ts.T())
	ts.service = New(ts.cfg, ts.storage)
	order := oService.New(ts.storage)
	ts.service.SetOrder(order)
}

func (ts *UserServiceTestSuite) TearDownTest() {
	// Reset mock
	for _, mc := range ts.storage.ExpectedCalls {
		mc.Unset()
	}
}

func (ts *UserServiceTestSuite) TearDownSubTest() {
	// Reset mock
	for _, mc := range ts.storage.ExpectedCalls {
		mc.Unset()
	}
}

func (ts *UserServiceTestSuite) setStorage(o store.Orderer, u store.Userer, l store.Loyaltier) {
	if o != nil {
		ts.storage.EXPECT().Order().Return(o)
	}

	if u != nil {
		ts.storage.EXPECT().User().Return(u)
	}

	if l != nil {
		ts.storage.EXPECT().Loyalty().Return(l)
	}
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
