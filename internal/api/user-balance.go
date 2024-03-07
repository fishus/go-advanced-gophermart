package api

import (
	"encoding/json"
	"net/http"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

// userBalance Получение баланса пользователя
func (s *server) userBalance(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Аутентификация пользователя
	token, err := s.auth(r)
	if err != nil {
		JSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	balance, err := s.service.Loyalty().UserBalance(r.Context(), token.UserID)
	if err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		logger.Log.Error(err.Error())
		return
	}

	type LoyaltyBalanceResult struct {
		Current   float64 `json:"current"`   // Текущий баланс
		Withdrawn float64 `json:"withdrawn"` // Списано за всё время
	}

	res := LoyaltyBalanceResult{
		Current: func() float64 {
			f, _ := balance.Current.Float64()
			return f
		}(),
		Withdrawn: func() float64 {
			f, _ := balance.Withdrawn.Float64()
			return f
		}(),
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Log.Debug(err.Error(), logger.Any("data", res))
		return
	}
}
