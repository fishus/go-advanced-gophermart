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

type Servicer interface {
	User() service.Userer
	Order() service.Orderer
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
	isShutdown atomic.Bool
}

func NewAccrual(cfg *Config, service Servicer) *daemon {
	client := resty.New().
		SetBaseURL("http://" + cfg.APIAddr).
		SetTimeout(cfg.RequestTimeout)

	if cfg.WorkersNum <= 0 {
		cfg.WorkersNum = 1
	}

	d := &daemon{
		cfg:      cfg,
		client:   client,
		service:  service,
		chOrders: make(chan models.Order),
		wg:       &sync.WaitGroup{},
	}
	d.delayCond = sync.NewCond(&sync.Mutex{})
	return d
}

func (d *daemon) Run(ctx context.Context) error {
	d.runOnce.Do(func() {
		logger.Log.Info("Running accrual daemon")

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
	d.isShutdown.Store(true)
	d.wg.Wait()
	close(d.chOrders)
	return nil
}

var _ io.Closer = (*daemon)(nil)
