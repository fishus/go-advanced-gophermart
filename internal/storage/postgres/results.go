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

type LoyaltyBalanceResult struct {
	UserID    models.UserID `db:"user_id"`   // ID пользователя
	Current   float64       `db:"current"`   // Текущий баланс
	Accrued   float64       `db:"accrued"`   // Начислено за всё время
	Withdrawn float64       `db:"withdrawn"` // Списано за всё время
}

type LoyaltyHistoryResult struct {
	UserID      models.UserID `db:"user_id"`      // ID пользователя
	OrderNum    string        `db:"order_num"`    // Номер заказа
	Accrual     float64       `db:"accrual"`      // Начисление
	Withdrawal  float64       `db:"withdrawal"`   // Списание
	ProcessedAt time.Time     `db:"processed_at"` // Дата зачисления или списания
}

func listResultsToLoyaltyHistory(results []LoyaltyHistoryResult) []models.LoyaltyHistory {
	history := make([]models.LoyaltyHistory, 0)
	for _, res := range results {
		h := models.LoyaltyHistory(res)
		history = append(history, h)
	}
	return history
}
