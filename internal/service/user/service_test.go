package user

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
