package api

import (
	"errors"
	"io"
	"net/http"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
)

// orderAdd Загрузка номера заказа для расчёта
func (s *server) orderAdd(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Аутентификация пользователя
	token, err := s.auth(r)
	if err != nil {
		JSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	orderNum, err := io.ReadAll(r.Body)
	if err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderID, err := s.service.Order().Add(r.Context(), token.UserID, string(orderNum))
	if err != nil {
		var validErr *serviceErr.ValidationError
		if errors.As(err, &validErr) {
			JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		// номер заказа уже был загружен этим пользователем
		if errors.Is(err, serviceErr.ErrOrderAlreadyExists) {
			JSONError(w, serviceErr.ErrOrderAlreadyExists.Error(), http.StatusOK)
			return
		}
		if errors.Is(err, serviceErr.ErrIncorrectData) {
			JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		// номер заказа уже был загружен другим пользователем
		if errors.Is(err, serviceErr.ErrOrderWrongOwner) {
			JSONError(w, serviceErr.ErrOrderWrongOwner.Error(), http.StatusConflict)
			return
		}
		// неверный формат номера заказа
		if errors.Is(err, serviceErr.ErrOrderWrongNum) {
			JSONError(w, serviceErr.ErrOrderWrongNum.Error(), http.StatusUnprocessableEntity)
			return
		}
		JSONError(w, err.Error(), http.StatusInternalServerError)
		logger.Log.Error(err.Error())
		return
	}
	logger.Log.Info("Registered new order",
		logger.String("OrderID", orderID.String()),
		logger.String("orderNum", string(orderNum)),
	)

	// новый номер заказа принят в обработку
	w.WriteHeader(http.StatusAccepted)
}
