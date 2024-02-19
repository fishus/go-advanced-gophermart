package order

import (
	"context"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (s *service) ResetProcessingStatus(ctx context.Context) error {
	return s.storage.OrderResetProcessingStatus(ctx)
}

// MoveToProcessing selects N new orders and changes their status to "processing" and then returns a list of these orders
func (s *service) MoveToProcessing(ctx context.Context, limit int) ([]models.Order, error) {
	//list, err = s.storage.OrderMoveToProcessing(ctx, limit)

	// List of orders before status changes
	list, err := s.storage.OrdersByFilter(ctx, limit, store.WithOrderStatus(models.OrderStatusNew), store.WithOrderBy(store.OrderByUploadedAt, store.OrderByAsc))
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	var idList []models.OrderID
	for _, order := range list {
		idList = append(idList, order.ID)
	}

	err = s.storage.OrderSetStatus(ctx, idList, models.OrderStatusProcessing)
	if err != nil {
		return nil, err
	}

	// List of orders after status changes
	list, err = s.storage.OrdersByFilter(ctx, limit, store.WithOrderIDList(idList...), store.WithOrderBy(store.OrderByUploadedAt, store.OrderByAsc))
	if err != nil {
		return nil, err
	}

	return list, nil
}
