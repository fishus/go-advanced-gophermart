package accrual

import (
	"context"
	"io"

	"github.com/go-resty/resty/v2"

	"github.com/fishus/go-advanced-gophermart/internal/service"
)

type Servicer interface {
	User() service.Userer
	Order() service.Orderer
}

type daemon struct {
	cfg     *Config
	client  *resty.Client
	service Servicer
}

func NewAccrual(cfg *Config, service Servicer) *daemon {
	client := resty.New().SetBaseURL("http://" + cfg.APIAddr)
	d := &daemon{
		cfg:     cfg,
		client:  client,
		service: service,
	}
	return d
}

func (d *daemon) Run(ctx context.Context) error {
	chNewOrders := d.PushNewOrders(ctx, d.cfg.LimitNewOrders)
	chDelayedOrders := d.workerDelayed(ctx, chNewOrders)
	d.workerGetOrderAccrual(ctx, chNewOrders, chDelayedOrders)
	return nil
}

func (d *daemon) Close() error {
	return nil
}

var _ io.Closer = (*daemon)(nil)
