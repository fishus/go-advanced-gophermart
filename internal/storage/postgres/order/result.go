package order

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

// OrderResult Заказ
type OrderResult struct {
	ID         models.OrderID     `db:"id"`          // ID заказа
	UserID     models.UserID      `db:"user_id"`     // ID пользователя
	Num        string             `db:"num"`         // Номер заказа
	Accrual    decimal.Decimal    `db:"accrual"`     // Начислено баллов лояльности
	Status     models.OrderStatus `db:"status"`      // Статус заказа
	UploadedAt time.Time          `db:"uploaded_at"` // Дата и время добавления заказа
	UpdatedAt  time.Time          `db:"updated_at"`  // Дата и время обновления статуса заказа
}

func listResultsToOrders(results []OrderResult) []models.Order {
	orders := make([]models.Order, 0)
	for _, res := range results {
		order := models.Order(res)
		orders = append(orders, order)
	}
	return orders
}
