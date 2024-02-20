package accrual

import (
	"errors"
	"time"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
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
	case OrderAccrualStatusRegistered:
	case OrderAccrualStatusInvalid:
	case OrderAccrualStatusProcessing:
	case OrderAccrualStatusProcessed:
		return nil
	case OrderAccrualStatusUndefined:
		return errors.New("order accrual status is not defined")
	}
	return errors.New("incorrect order accrual status")
}

func (s OrderAccrualStatus) String() string {
	return string(s)
}

// OrderAccrual Hасчёт начисления вознаграждения за заказ
type OrderAccrual struct {
	Num     string             // Номер заказа
	Status  OrderAccrualStatus // Статус расчёта начисления вознаграждения за заказ
	Accrual float64            // Начислено баллов лояльности // FIXME
}

type delayedOrder struct {
	order models.Order
	delay time.Time
}
