package order

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	apiCommon "github.com/fishus/go-advanced-gophermart/internal/api/common"
	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

// List Список загруженных номеров заказов
func (a *api) List(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Аутентификация пользователя
	token, err := a.auth(r)
	if err != nil {
		apiCommon.JSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	list, err := a.service.Order().ListByUser(r.Context(), token.UserID)
	if err != nil {
		apiCommon.JSONError(w, err.Error(), http.StatusInternalServerError)
		logger.Log.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if len(list) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	type OrderResult struct {
		Num        string             `json:"number"`            // Номер заказа
		Accrual    float64            `json:"accrual,omitempty"` // Начислено баллов лояльности
		Status     models.OrderStatus `json:"status"`            // Статус заказа
		UploadedAt time.Time          `json:"uploaded_at"`       // Дата и время добавления заказа
	}

	ordersList := make([]OrderResult, 0)
	for _, order := range list {
		o := OrderResult{
			Num:        order.Num,
			Accrual:    order.Accrual.InexactFloat64(),
			Status:     order.Status,
			UploadedAt: order.UploadedAt,
		}
		ordersList = append(ordersList, o)
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(ordersList); err != nil {
		logger.Log.Debug(err.Error(), logger.Any("data", ordersList))
		return
	}
}
