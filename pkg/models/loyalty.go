package models

import "time"

// LoyaltyBalance Баланс пользователя в программе лояльности
type LoyaltyBalance struct {
	UserID    UserID  // ID пользователя
	Current   float64 // Текущий баланс
	Accrued   float64 // Начислено за всё время
	Withdrawn float64 // Списано за всё время
}

// LoyaltyHistory История изменения баланса пользователя в программе лояльности (начисления и списания баллов лояльности)
type LoyaltyHistory struct {
	UserID      UserID    // ID пользователя
	OrderID     OrderID   // ID заказа
	Accrual     float64   // Начисление
	Withdrawal  float64   // Списание
	ProcessedAt time.Time // Дата зачисления или списания
}
