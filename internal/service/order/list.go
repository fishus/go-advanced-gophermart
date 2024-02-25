package order

import (
	"context"
	"errors"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

// ListNew returns orders in "new" status
func (s *service) ListNew(ctx context.Context) ([]models.Order, error) {
	list, err := s.storage.OrdersByFilter(ctx, 0, store.WithOrderStatus(models.OrderStatusNew), store.WithOrderBy(store.OrderByUploadedAt, store.OrderByAsc))
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return list, nil
}

// ListProcessing returns orders in "processing" status
func (s *service) ListProcessing(ctx context.Context, limit int) ([]models.Order, error) {
	list, err := s.storage.OrdersByFilter(ctx, limit, store.WithOrderStatus(models.OrderStatusProcessing), store.WithOrderBy(store.OrderByUpdatedAt, store.OrderByAsc))
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return list, nil
}
