package order

import (
	"context"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (s *service) UpdateStatus(ctx context.Context, id models.OrderID, status models.OrderStatus) error {
	if err := status.Validate(); err != nil {
		return err
	}
	return s.storage.Order().UpdateStatus(ctx, id, status)
}
