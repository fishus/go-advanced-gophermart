package app

import (
	"testing"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

func init() {
	if testing.Testing() {
		return
	}
	Config = initConfig()
	_ = logger.Initialize(Config.LogLevel())
}
