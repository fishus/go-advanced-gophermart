package accrual

import (
	"context"
	"time"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

// addNewOrders Add to queue all new orders (in case the service has been restarted)
func (d *daemon) addNewOrders(ctx context.Context) {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()

		list, err := d.service.Order().ListNew(ctx)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		if len(list) == 0 {
			return
		}

		for _, order := range list {
			d.chOrders <- order
		}
	}()
}

// addProcessingOrders Add to queue "processing" orders
func (d *daemon) addProcessingOrders(ctx context.Context) {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-d.chShutdown:
				return
			case <-ticker.C:
				list, err := d.service.Order().ListProcessing(ctx, d.cfg.WorkersNum)
				if err != nil {
					logger.Log.Error(err.Error())
					return
				}

				if len(list) == 0 {
					continue
				}

				for _, order := range list {
					d.chOrders <- order
				}
			}
		}
	}()
}

// AddNewOrder Add a new order to the processing pipeline
func (d *daemon) AddNewOrder(ctx context.Context, order models.Order) {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		d.chOrders <- order
	}()
}
