package models

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
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
	case OrderStatusNew,
		OrderStatusProcessing,
		OrderStatusInvalid,
		OrderStatusProcessed:
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
	ID         OrderID         // ID заказа
	UserID     UserID          // ID пользователя
	Num        string          // Номер заказа
	Accrual    decimal.Decimal // Начислено баллов лояльности
	Status     OrderStatus     // Статус заказа
	UploadedAt time.Time       // Дата и время добавления заказа
	UpdatedAt  time.Time       // Дата и время обновления статуса заказа
}
