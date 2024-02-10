package main

import (
	"context"

	"github.com/fishus/go-advanced-gophermart/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.Shutdown(cancel)

	<-ctx.Done()
}
