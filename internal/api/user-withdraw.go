package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
)

// userWithdraw Запрос на списание средств
func (s *server) userWithdraw(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Аутентификация пользователя
	token, err := s.auth(r)
	if err != nil {
		JSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	type reqData struct {
		Num string  `json:"order"` // Номер заказа
		Sum float64 `json:"sum"`   // Сумма баллов к списанию в счёт оплаты
	}

	var data reqData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		logger.Log.Debug(err.Error())
		return
	}

	err = s.service.User().LoyaltyAddWithdraw(r.Context(), token.UserID, data.Num, data.Sum)
	if err != nil {
		retCode := http.StatusInternalServerError

		switch {
		case errors.Is(err, serviceErr.ErrOrderWrongNum):
			retCode = http.StatusUnprocessableEntity

		case errors.Is(err, serviceErr.ErrIncorrectData):
			retCode = http.StatusBadRequest

		case errors.Is(err, serviceErr.ErrLowBalance):
			retCode = http.StatusPaymentRequired

		default:
			logger.Log.Error(err.Error())
		}

		JSONError(w, err.Error(), retCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}
