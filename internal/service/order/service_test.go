package order

import (
	"testing"

	"github.com/stretchr/testify/suite"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	storageMocks "github.com/fishus/go-advanced-gophermart/internal/storage/mocks"
)

type OrderServiceTestSuite struct {
	suite.Suite
	storage *storageMocks.Storager
	*service
}

func (ts *OrderServiceTestSuite) SetupSuite() {
	ts.storage = storageMocks.NewStorager(ts.T())
	ts.service = New(ts.storage)
}

func (ts *OrderServiceTestSuite) TearDownTest() {
	// Reset mock
	for _, mc := range ts.storage.ExpectedCalls {
		mc.Unset()
	}
}

func (ts *OrderServiceTestSuite) TearDownSubTest() {
	// Reset mock
	for _, mc := range ts.storage.ExpectedCalls {
		mc.Unset()
	}
}

func (ts *OrderServiceTestSuite) setStorage(o store.Orderer, u store.Userer, l store.Loyaltier) {
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

func TestOrderService(t *testing.T) {
	suite.Run(t, new(OrderServiceTestSuite))
}
