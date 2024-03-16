package accrual

import (
	"errors"

	"github.com/shopspring/decimal"
)

// OrderAccrualStatus Статус расчёта начисления вознаграждения за заказ
type OrderAccrualStatus string

const (
	OrderAccrualStatusUndefined  OrderAccrualStatus = ""
	OrderAccrualStatusRegistered OrderAccrualStatus = "REGISTERED" // заказ зарегистрирован, но вознаграждение не рассчитано;
	OrderAccrualStatusInvalid    OrderAccrualStatus = "INVALID"    // заказ не принят к расчёту, и вознаграждение не будет начислено;
	OrderAccrualStatusProcessing OrderAccrualStatus = "PROCESSING" // расчёт начисления в процессе;
	OrderAccrualStatusProcessed  OrderAccrualStatus = "PROCESSED"  // расчёт начисления окончен;
)

func (s OrderAccrualStatus) Validate() error {
	switch s {
	case OrderAccrualStatusRegistered,
		OrderAccrualStatusInvalid,
		OrderAccrualStatusProcessing,
		OrderAccrualStatusProcessed:
		return nil
	case OrderAccrualStatusUndefined:
		return errors.New("order accrual status is not defined")
	}
	return errors.New("incorrect order accrual status")
}

func (s OrderAccrualStatus) String() string {
	return string(s)
}

// OrderAccrual Расчёт начисления вознаграждения за заказ
type OrderAccrual struct {
	Num     string             `json:"order"`   // Номер заказа
	Status  OrderAccrualStatus `json:"status"`  // Статус расчёта начисления вознаграждения за заказ
	Accrual decimal.Decimal    `json:"accrual"` // Начислено баллов лояльности
}
