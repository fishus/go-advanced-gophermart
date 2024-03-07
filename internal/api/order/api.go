package order

import (
	"context"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/service"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=Servicer --with-expecter
type Servicer interface {
	User() service.Userer
	Order() service.Orderer
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=AccrualDaemon  --with-expecter
type AccrualDaemon interface {
	AddNewOrder(context.Context, models.Order)
}

type api struct {
	service Servicer
	daemon  AccrualDaemon
}

func NewAPI(service Servicer, daemon AccrualDaemon) *api {
	return &api{
		service: service,
		daemon:  daemon,
	}
}
