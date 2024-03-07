package loyalty

import (
	"testing"

	"github.com/stretchr/testify/suite"

	oService "github.com/fishus/go-advanced-gophermart/internal/service/order"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	storageMocks "github.com/fishus/go-advanced-gophermart/internal/storage/mocks"
)

type LoyaltyServiceTestSuite struct {
	suite.Suite
	storage *storageMocks.Storager
	*service
}

func (ts *LoyaltyServiceTestSuite) SetupSuite() {
	ts.storage = storageMocks.NewStorager(ts.T())
	ts.service = New(ts.storage)
	order := oService.New(ts.storage)
	ts.service.SetOrder(order)
}

func (ts *LoyaltyServiceTestSuite) TearDownTest() {
	// Reset mock
	for _, mc := range ts.storage.ExpectedCalls {
		mc.Unset()
	}
}

func (ts *LoyaltyServiceTestSuite) TearDownSubTest() {
	// Reset mock
	for _, mc := range ts.storage.ExpectedCalls {
		mc.Unset()
	}
}

func (ts *LoyaltyServiceTestSuite) setStorage(o store.Orderer, u store.Userer, l store.Loyaltier) {
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

func TestLoyaltyService(t *testing.T) {
	suite.Run(t, new(LoyaltyServiceTestSuite))
}
