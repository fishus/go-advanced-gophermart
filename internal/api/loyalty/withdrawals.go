package loyalty

import (
	"encoding/json"
	"net/http"
	"time"

	apiCommon "github.com/fishus/go-advanced-gophermart/internal/api/common"
	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

// Withdrawals Информации о выводе средств
func (a *api) Withdrawals(w http.ResponseWriter, r *http.Request) {
	// Аутентификация пользователя
	token, err := a.auth(r)
	if err != nil {
		apiCommon.JSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Факты выводов в выдаче должны быть отсортированы по времени вывода от самых старых к самым новым.
	// Формат даты — RFC3339.

	history, err := a.service.Loyalty().UserWithdrawals(r.Context(), token.UserID)
	if err != nil {
		apiCommon.JSONError(w, err.Error(), http.StatusInternalServerError)
		logger.Log.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if len(history) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	type LoyaltyHistoryResult struct {
		OrderNum    string    `json:"order"`        // Номер заказа
		Withdrawal  float64   `json:"sum"`          // Списание
		ProcessedAt time.Time `json:"processed_at"` // Дата зачисления или списания
	}

	historyList := make([]LoyaltyHistoryResult, 0)
	for _, h := range history {
		o := LoyaltyHistoryResult{
			OrderNum:    h.OrderNum,
			Withdrawal:  h.Withdrawal.InexactFloat64(),
			ProcessedAt: h.ProcessedAt,
		}
		historyList = append(historyList, o)
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(historyList); err != nil {
		logger.Log.Debug(err.Error(), logger.Any("data", historyList))
		return
	}
}
