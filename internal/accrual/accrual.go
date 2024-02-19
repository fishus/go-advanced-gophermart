package accrual

import (
	"context"
	"github.com/fishus/go-advanced-gophermart/internal/logger"
	"github.com/fishus/go-advanced-gophermart/internal/service"
	"io"
)

type Servicer interface {
	User() service.Userer
	Order() service.Orderer
}

type daemon struct {
	cfg     *Config
	service Servicer
}

func NewAccrual(cfg *Config, service Servicer) *daemon {
	d := &daemon{
		cfg:     cfg,
		service: service,
	}
	return d
}

func (d *daemon) Run(ctx context.Context) error {
	chNewOrders := d.PushNewOrders(ctx, d.cfg.LimitNewOrders)
	go func() {
		for order := range chNewOrders {
			select {
			case <-ctx.Done():
				return
			default:
			}
			logger.Log.Info("New order", logger.Any("order", order))
		}
	}()
	// TODO Run workers
	return nil
}

func (d *daemon) Close() error {
	// TODO
	//logger.Log.Info("Shutdown api server")
	//ctx := context.Background()
	//return s.server.Shutdown(ctx)
	//close(s.chOrder)
	return nil
}

var _ io.Closer = (*daemon)(nil)
