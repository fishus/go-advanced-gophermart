package order

import (
	"testing"

	"github.com/stretchr/testify/suite"

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

func TestOrderService(t *testing.T) {
	suite.Run(t, new(OrderServiceTestSuite))
}
