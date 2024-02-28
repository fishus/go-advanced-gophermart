package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	oService "github.com/fishus/go-advanced-gophermart/internal/service/order"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

type UserServiceTestSuite struct {
	suite.Suite
	cfg     *Config
	storage *store.MockStorage
	*service
}

func (ts *UserServiceTestSuite) SetupSuite() {
	ts.cfg = &Config{
		JWTExpires:   15 * time.Minute,
		JWTSecretKey: "TestSecretKey",
	}
	ts.storage = new(store.MockStorage)
	ts.service = New(ts.cfg, ts.storage)
	order := oService.New(ts.storage)
	ts.service.SetOrder(order)
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
