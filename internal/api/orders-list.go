package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

// ordersList Список загруженных номеров заказов
func (s *server) ordersList(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Аутентификация пользователя
	token, err := s.auth(r)
	if err != nil {
		JSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	list, err := s.service.Order().ListByUser(r.Context(), token.UserID)
	if err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		logger.Log.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if len(list) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	type OrderResult struct {
		ID         models.OrderID     `json:"-"`                 // ID заказа
		UserID     models.UserID      `json:"-"`                 // ID пользователя
		Num        string             `json:"number"`            // Номер заказа
		Accrual    float64            `json:"accrual,omitempty"` // Начислено баллов лояльности // FIXME
		Status     models.OrderStatus `json:"status"`            // Статус заказа
		UploadedAt time.Time          `json:"uploaded_at"`       // Дата и время добавления заказа
		UpdatedAt  time.Time          `json:"-"`                 // Дата и время обновления статуса заказа
	}

	ordersList := make([]OrderResult, 0)
	for _, order := range list {
		o := OrderResult(order)
		ordersList = append(ordersList, o)
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(ordersList); err != nil {
		logger.Log.Debug(err.Error(), logger.Any("data", ordersList))
		return
	}
}
