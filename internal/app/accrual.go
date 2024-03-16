package app

import (
	"context"
	"time"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/accrual"
	"github.com/fishus/go-advanced-gophermart/internal/service"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

type AccrualDaemon interface {
	AddNewOrder(context.Context, models.Order)
}

func RunAccrualWorkers(ctx context.Context, storage store.Storager) (AccrualDaemon, error) {
	serv := service.New(&service.Config{}, storage)

	accConfig := &accrual.Config{
		APIAddr:        Config.AccrualAddr(),
		RequestTimeout: 5 * time.Second,
		MaxAttempts:    3,
		WorkersNum:     3,
	}

	acc := accrual.NewAccrual(accConfig, serv)
	Closers = append(Closers, acc)
	err := acc.Run(ctx)
	return acc, err
}
