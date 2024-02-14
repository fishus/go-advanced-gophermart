package order

import (
	"context"
	"errors"
	"strings"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
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

	// Проверка, был ли такой заказ загружен ранее
	o, err := s.storage.OrderByFilter(ctx, store.WithOrderNum(orderNum))
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return orderID, err
	}
	if o != nil {
		// номер заказа уже был загружен другим пользователем
		if o.UserID == userID {
			return orderID, serviceErr.ErrOrderAlreadyExists
		}
		// номер заказа уже был загружен этим пользователем
		return orderID, serviceErr.ErrOrderWrongOwner
	}

	orderID, err = s.storage.OrderAdd(ctx, order)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			err = serviceErr.ErrOrderAlreadyExists
			return
		}
		if errors.Is(err, store.ErrIncorrectData) {
			err = serviceErr.ErrIncorrectData
			return
		}
		return orderID, err
	}

	return orderID, nil
}
