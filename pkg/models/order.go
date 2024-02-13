package models

import (
	"errors"
	"time"
)

type OrderID string

func (id OrderID) String() string {
	return string(id)
}

// OrderStatus Статус заказа
type OrderStatus string

const (
	OrderStatusUndefined  OrderStatus = ""
	OrderStatusNew        OrderStatus = "NEW"        // Заказ загружен в систему, но не попал в обработку;
	OrderStatusProcessing OrderStatus = "PROCESSING" // Вознаграждение за заказ рассчитывается;
	OrderStatusInvalid    OrderStatus = "INVALID"    // Система расчёта вознаграждений отказала в расчёте и вознаграждение не будет начислено;
	OrderStatusProcessed  OrderStatus = "PROCESSED"  // Данные по заказу проверены и информация о расчёте успешно получена.
)

func (s OrderStatus) Validate() error {
	switch s {
	case OrderStatusNew:
	case OrderStatusProcessing:
	case OrderStatusInvalid:
	case OrderStatusProcessed:
		return nil
	case OrderStatusUndefined:
		return errors.New("order status is not defined")
	}
	return errors.New("incorrect order status")
}

func (s OrderStatus) String() string {
	return string(s)
}

// Order Заказ
type Order struct {
	ID         OrderID     `db:"id"`                          // ID заказа
	UserID     UserID      `db:"user_id" validate:"required"` // ID пользователя
	Num        string      `db:"num" validate:"required"`     // Номер заказа
	Accrual    float64     `db:"accrual"`                     // Начислено баллов лояльности
	Status     OrderStatus `db:"status"`                      // Статус заказа
	UploadedAt time.Time   `db:"uploaded_at"`                 // Дата и время добавления заказа
}
