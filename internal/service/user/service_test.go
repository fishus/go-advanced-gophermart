package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	
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
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
