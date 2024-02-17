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

	// Проверка, был ли такой заказ загружен ранее
	o, err := s.storage.OrderByFilter(ctx, store.WithOrderNum(orderNum))
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return
	}
	if err == nil {
		// номер заказа уже был загружен другим пользователем
		if o.UserID == userID {
			err = serviceErr.ErrOrderAlreadyExists
			return
		}
		// номер заказа уже был загружен этим пользователем
		err = serviceErr.ErrOrderWrongOwner
		return
	}

	// TODO Рассмотр. возможность объединения в один запрос с проверкой
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
	}
	return
}
