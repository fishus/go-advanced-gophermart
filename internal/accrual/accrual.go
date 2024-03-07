package accrual

import (
	"context"
	"io"
	"sync"
	"sync/atomic"

	"github.com/go-resty/resty/v2"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
	"github.com/fishus/go-advanced-gophermart/internal/service"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=Servicer --with-expecter
type Servicer interface {
	User() service.Userer
	Order() service.Orderer
	Loyalty() service.Loyaltier
}

type daemon struct {
	cfg        *Config
	client     *resty.Client
	service    Servicer
	chOrders   chan models.Order
	runOnce    sync.Once
	wg         *sync.WaitGroup
	delay      atomic.Bool
	delayCond  *sync.Cond
	chShutdown chan struct{}
}

func NewAccrual(cfg *Config, service Servicer) *daemon {
	client := resty.New().
		SetBaseURL(cfg.APIAddr).
		SetTimeout(cfg.RequestTimeout)

	if cfg.WorkersNum <= 0 {
		cfg.WorkersNum = 1
	}

	d := &daemon{
		cfg:        cfg,
		client:     client,
		service:    service,
		chOrders:   make(chan models.Order),
		wg:         &sync.WaitGroup{},
		chShutdown: make(chan struct{}),
	}
	d.delayCond = sync.NewCond(&sync.Mutex{})
	return d
}

func (d *daemon) Run(ctx context.Context) error {
	d.runOnce.Do(func() {
		logger.Log.Info("Running accrual daemon on: " + d.cfg.APIAddr)

		d.addNewOrders(ctx)
		d.addProcessingOrders(ctx)
		for i := 0; i < d.cfg.WorkersNum; i++ {
			d.workerRequestOrderAccrual(ctx)
		}
	})
	return nil
}

func (d *daemon) Close() error {
	logger.Log.Info("Shutdown accrual daemon")
	close(d.chShutdown)
	d.wg.Wait()
	close(d.chOrders)
	return nil
}

var _ io.Closer = (*daemon)(nil)
