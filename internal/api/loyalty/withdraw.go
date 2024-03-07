package loyalty

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/shopspring/decimal"

	apiCommon "github.com/fishus/go-advanced-gophermart/internal/api/common"
	"github.com/fishus/go-advanced-gophermart/internal/app/config"
	"github.com/fishus/go-advanced-gophermart/internal/logger"
	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
)

// Withdraw Запрос на списание средств
func (a *api) Withdraw(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Аутентификация пользователя
	token, err := a.auth(r)
	if err != nil {
		apiCommon.JSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	type reqData struct {
		Num string          `json:"order"` // Номер заказа
		Sum decimal.Decimal `json:"sum"`   // Сумма баллов к списанию в счёт оплаты
	}

	var data reqData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		apiCommon.JSONError(w, err.Error(), http.StatusBadRequest)
		logger.Log.Debug(err.Error())
		return
	}
	data.Sum = data.Sum.Round(config.DecimalExponent)

	err = a.service.Loyalty().AddWithdraw(r.Context(), token.UserID, data.Num, data.Sum)
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

		apiCommon.JSONError(w, err.Error(), retCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}
