package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// LoyaltyBalance Баланс пользователя в программе лояльности
type LoyaltyBalance struct {
	UserID    UserID          // ID пользователя
	Current   decimal.Decimal // Текущий баланс
	Accrued   decimal.Decimal // Начислено за всё время
	Withdrawn decimal.Decimal // Списано за всё время
}

// LoyaltyHistory История изменения баланса пользователя в программе лояльности (начисления и списания баллов лояльности)
type LoyaltyHistory struct {
	UserID      UserID          // ID пользователя
	OrderNum    string          // Номер заказа
	Accrual     decimal.Decimal // Начисление
	Withdrawal  decimal.Decimal // Списание
	ProcessedAt time.Time       // Дата зачисления или списания
}
