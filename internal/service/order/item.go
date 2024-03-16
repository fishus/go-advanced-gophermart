package order

import (
	"context"
	"errors"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (s *service) GetByID(ctx context.Context, id models.OrderID) (order models.Order, err error) {
	order, err = s.storage.Order().GetByID(ctx, id)
	if err != nil && errors.Is(err, store.ErrNotFound) {
		err = serviceErr.ErrOrderNotFound
	}
	return
}
