package accrual

import (
	"context"
	"time"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

func (d *daemon) PushNewOrders(ctx context.Context, limit int) chan models.Order {
	ch := make(chan models.Order, limit)

	go func() {
		select {
		case <-ctx.Done():
			close(ch)
			return
		default:
		}

		// на первой итерации сбрасываем статус всех незавершенных заказов (на случай, если сервис перезапустился)
		err := d.service.Order().ResetProcessingStatus(ctx)
		if err != nil {
			logger.Log.Error(err.Error())
		}

		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			case <-ticker.C:
				list, err := d.service.Order().MoveToProcessing(ctx, d.cfg.LimitNewOrders)
				if err != nil {
					logger.Log.Error(err.Error())
				}
				for _, order := range list {
					ch <- order
				}
			}
		}
	}()

	return ch
}
