package order

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type OrderServiceTestSuite struct {
	suite.Suite
}

func TestOrderService(t *testing.T) {
	suite.Run(t, new(OrderServiceTestSuite))
}
