package postgres

import (
	"time"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

// UserResult Пользователь
type UserResult struct {
	ID        models.UserID `db:"id"`         // ID пользователя
	Username  string        `db:"username"`   // Логин
	Password  string        `db:"-"`          // Пароль
	CreatedAt time.Time     `db:"created_at"` // Дата регистрации
}

// OrderResult Заказ
type OrderResult struct {
	ID         models.OrderID     `db:"id"`          // ID заказа
	UserID     models.UserID      `db:"user_id"`     // ID пользователя
	Num        string             `db:"num"`         // Номер заказа
	Accrual    float64            `db:"accrual"`     // Начислено баллов лояльности // FIXME Хранить в int в копейках
	Status     models.OrderStatus `db:"status"`      // Статус заказа
	UploadedAt time.Time          `db:"uploaded_at"` // Дата и время добавления заказа
}

func listResultsToOrders(results []OrderResult) []models.Order {
	orders := make([]models.Order, 0)
	for _, res := range results {
		order := models.Order(res)
		orders = append(orders, order)
	}
	return orders
}
