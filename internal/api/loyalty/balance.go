package loyalty

import (
	"encoding/json"
	"net/http"

	apiCommon "github.com/fishus/go-advanced-gophermart/internal/api/common"
	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

// Balance Получение баланса пользователя
func (a *api) Balance(w http.ResponseWriter, r *http.Request) {
	// Аутентификация пользователя
	token, err := a.auth(r)
	if err != nil {
		apiCommon.JSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	balance, err := a.service.Loyalty().UserBalance(r.Context(), token.UserID)
	if err != nil {
		apiCommon.JSONError(w, err.Error(), http.StatusInternalServerError)
		logger.Log.Error(err.Error())
		return
	}

	type LoyaltyBalanceResult struct {
		Current   float64 `json:"current"`   // Текущий баланс
		Withdrawn float64 `json:"withdrawn"` // Списано за всё время
	}

	res := LoyaltyBalanceResult{
		Current:   balance.Current.InexactFloat64(),
		Withdrawn: balance.Withdrawn.InexactFloat64(),
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Log.Debug(err.Error(), logger.Any("data", res))
		return
	}
}
