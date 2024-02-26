package api

import (
	"encoding/json"
	"net/http"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

// userBalance Получение баланса пользователя
func (s *server) userBalance(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Аутентификация пользователя
	token, err := s.auth(r)
	if err != nil {
		JSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	balance, err := s.service.User().LoyaltyUserBalance(r.Context(), token.UserID)
	if err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		logger.Log.Error(err.Error())
		return
	}

	type LoyaltyBalanceResult struct {
		UserID    models.UserID `json:"-"`         // ID пользователя
		Current   float64       `json:"current"`   // Текущий баланс
		Accrued   float64       `json:"-"`         // Начислено за всё время
		Withdrawn float64       `json:"withdrawn"` // Списано за всё время
	}

	res := LoyaltyBalanceResult(balance)

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Log.Debug(err.Error(), logger.Any("data", res))
		return
	}
}
