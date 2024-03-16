package app

import (
	"testing"

	"github.com/fishus/go-advanced-gophermart/internal/app/config"
	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

var Config config.Config

func init() {
	if testing.Testing() {
		return
	}
	Config = config.InitConfig()
	_ = logger.Initialize(Config.LogLevel())
}
