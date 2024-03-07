package loyalty

import (
	"github.com/fishus/go-advanced-gophermart/internal/service"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=Servicer --with-expecter
type Servicer interface {
	User() service.Userer
	Loyalty() service.Loyaltier
}

type api struct {
	service Servicer
}

func NewAPI(service Servicer) *api {
	return &api{
		service: service,
	}
}
