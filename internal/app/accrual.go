package app

import (
	"context"

	"github.com/fishus/go-advanced-gophermart/internal/accrual"
	"github.com/fishus/go-advanced-gophermart/internal/service"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func RunAccrualWorkers(ctx context.Context, storage store.Storager) error {
	serv := service.New(&service.Config{}, storage)

	accConfig := &accrual.Config{
		APIAddr: Config.AccrualAddr(),
		//LimitNewOrders: 100,
		LimitNewOrders: 3,
	}

	acc := accrual.NewAccrual(accConfig, serv)
	Closers = append(Closers, acc)
	err := acc.Run(ctx)
	return err
}
