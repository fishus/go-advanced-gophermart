package order

import (
	"testing"

	"github.com/stretchr/testify/suite"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

type OrderServiceTestSuite struct {
	suite.Suite
	storage *store.MockStorage
	*service
}

func (ts *OrderServiceTestSuite) SetupSuite() {
	ts.storage = new(store.MockStorage)
	ts.service = New(ts.storage)
}

func TestOrderService(t *testing.T) {
	suite.Run(t, new(OrderServiceTestSuite))
}
