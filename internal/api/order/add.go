package order

import (
	"context"
	"errors"
	"io"
	"net/http"

	apiCommon "github.com/fishus/go-advanced-gophermart/internal/api/common"
	"github.com/fishus/go-advanced-gophermart/internal/logger"
	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
)

// Add Загрузка номера заказа для расчёта
func (a *api) Add(w http.ResponseWriter, r *http.Request) {
	// Аутентификация пользователя
	token, err := a.auth(r)
	if err != nil {
		apiCommon.JSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	orderNum, err := io.ReadAll(r.Body)
	if err != nil {
		apiCommon.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderID, err := a.service.Order().Add(r.Context(), token.UserID, string(orderNum))
	if err != nil {
		var validErr *serviceErr.ValidationError
		if errors.As(err, &validErr) {
			apiCommon.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		// номер заказа уже был загружен этим пользователем
		if errors.Is(err, serviceErr.ErrOrderAlreadyExists) {
			apiCommon.JSONError(w, serviceErr.ErrOrderAlreadyExists.Error(), http.StatusOK)
			return
		}
		if errors.Is(err, serviceErr.ErrIncorrectData) {
			apiCommon.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		// номер заказа уже был загружен другим пользователем
		if errors.Is(err, serviceErr.ErrOrderWrongOwner) {
			apiCommon.JSONError(w, serviceErr.ErrOrderWrongOwner.Error(), http.StatusConflict)
			return
		}
		// неверный формат номера заказа
		if errors.Is(err, serviceErr.ErrOrderWrongNum) {
			apiCommon.JSONError(w, serviceErr.ErrOrderWrongNum.Error(), http.StatusUnprocessableEntity)
			return
		}
		apiCommon.JSONError(w, err.Error(), http.StatusInternalServerError)
		logger.Log.Error(err.Error())
		return
	}
	logger.Log.Info("Registered new order",
		logger.String("OrderID", orderID.String()),
		logger.String("orderNum", string(orderNum)),
	)

	ctx := context.Background()
	order, err := a.service.Order().GetByID(ctx, orderID)
	if err != nil {
		logger.Log.Error(err.Error())
		apiCommon.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.daemon.AddNewOrder(ctx, order)

	// новый номер заказа принят в обработку
	w.WriteHeader(http.StatusAccepted)
}
