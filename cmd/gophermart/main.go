package main

import (
	"context"

	"github.com/fishus/go-advanced-gophermart/internal/app"
	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.Shutdown(cancel)

	db, err := app.ConnDB(ctx)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	err = app.RunAccrualWorkers(ctx, db)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	err = app.RunAPIServer(ctx, db)
	if err != nil {
		logger.Log.Error(err.Error())
	}
}
