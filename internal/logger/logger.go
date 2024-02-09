package logger

import (
	"go.uber.org/zap"

	"github.com/fishus/go-advanced-gophermart/internal/config"
)

var Log *zap.Logger = zap.NewNop()

func init() {
	lvl, err := zap.ParseAtomicLevel(config.Config.LogLevel())
	if err != nil {
		return
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	l, err := cfg.Build()
	if err != nil {
		return
	}
	Log = l
}
