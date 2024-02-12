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

	err := app.RunAPIServer(ctx)
	logger.Log.Info(err.Error())
}
