package main

import (
	"context"

	"github.com/fishus/go-advanced-gophermart/internal/app"
	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	app.Shutdown(cancel)

	db, err := app.ConnDB(ctx)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	loyalty, err := app.RunAccrualWorkers(ctx, db)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	err = app.RunAPIServer(ctx, db, loyalty)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	<-ctx.Done()
}
