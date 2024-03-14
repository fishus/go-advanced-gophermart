package app

import (
	"context"
	"io"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

var (
	Closers      []io.Closer
	onceShutdown sync.Once
)

func CloseClosers() {
	slices.Reverse(Closers)

	for _, closer := range Closers {
		err := closer.Close()
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}
}

// Shutdown implements graceful app

func WaitSignalTerm(cancel context.CancelFunc) {
	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		Shutdown(cancel)
	}()
}

func Shutdown(cancel context.CancelFunc) {
	onceShutdown.Do(func() {
		CloseClosers()
		cancel()
	})
}
