package order

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
)

func (s *service) UpdateStatus(ctx context.Context, id models.OrderID, status models.OrderStatus) error {
	if err := status.Validate(); err != nil {
		return err
	}
	return s.storage.Order().UpdateStatus(ctx, id, status)
}

func (s *service) AddAccrual(ctx context.Context, id models.OrderID, accrual decimal.Decimal) error {
	if accrual.LessThan(decimal.NewFromFloat(0)) {
		return serviceErr.ErrIncorrectData
	}

	order, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if order.Status == models.OrderStatusProcessed {
		return serviceErr.ErrOrderRewardReceived
	}

	return s.storage.Order().AddAccrual(ctx, id, accrual)
}
