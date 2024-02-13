package order

import (
	"context"
	"errors"
	"strings"

	"github.com/gookit/validate"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

// Add Загрузка номера заказа для расчёта
func (s *service) Add(ctx context.Context, userID models.UserID, orderNum string) (models.OrderID, error) {
	var orderID models.OrderID

	orderNum = strings.TrimSpace(orderNum)

	order := &models.Order{
		UserID: userID,
		Num:    orderNum,
		Status: models.OrderStatusNew,
	}

	v := validate.Struct(order)
	if !v.Validate() {
		return orderID, serviceErr.NewValidationError(v.Errors)
	}

	// Проверка номера заказа на корректность с помощью алгоритма Луна
	if err := validateNumLuhn(orderNum); err != nil {
		return orderID, err
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

	orderID, err = s.storage.OrderAdd(ctx, *order)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			err = serviceErr.ErrOrderAlreadyExists
		}
		return orderID, err
	}

	return orderID, nil
}
