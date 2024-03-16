package order

import (
	"context"
	"errors"
	"strings"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

// Add Загрузка номера заказа для расчёта
func (s *service) Add(ctx context.Context, userID models.UserID, orderNum string) (orderID models.OrderID, err error) {
	orderNum = strings.TrimSpace(orderNum)

	order := models.Order{
		UserID: userID,
		Num:    orderNum,
		Status: models.OrderStatusNew,
	}

	err = validateOrder(order)
	if err != nil {
		return
	}

	// Проверка номера заказа на корректность с помощью алгоритма Луна
	if err = validateNumLuhn(orderNum); err != nil {
		return
	}

	orderID, err = s.storage.Order().Add(ctx, order)
	if err != nil {
		if errors.Is(err, store.ErrIncorrectData) {
			err = serviceErr.ErrIncorrectData
			return
		}
		// Если такой заказ загружен ранее
		o, oErr := s.storage.Order().GetByFilter(ctx, store.WithOrderNum(orderNum))
		if oErr == nil {
			// номер заказа уже был загружен этим пользователем
			if o.UserID == userID {
				err = serviceErr.ErrOrderAlreadyExists
				return
			}
			// номер заказа уже был загружен другим пользователем
			err = serviceErr.ErrOrderWrongOwner
			return
		}
	}
	return
}
