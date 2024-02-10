package main

import (
	"context"
	"fmt"

	"github.com/fishus/go-advanced-gophermart/internal/app"
	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.Shutdown(cancel)

	db, err := app.ConnDB(ctx)
	if err != nil {
		logger.Log.Panic(err.Error())
	}
	logger.Log.Info("DB storage", logger.String("db", fmt.Sprintf("%+v", db)))

	<-ctx.Done()
}
