package app

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

var Closers []io.Closer

func CloseClosers() {
	for _, closer := range Closers {
		if err := closer.Close(); err != nil {
			logger.Log.Error(err.Error())
		}
	}
}

// Shutdown implements graceful app
func Shutdown(cancel context.CancelFunc) {
	go func() {
		termSig := make(chan os.Signal, 1)
		signal.Notify(termSig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		<-termSig
		cancel()
		CloseClosers()
	}()
}
