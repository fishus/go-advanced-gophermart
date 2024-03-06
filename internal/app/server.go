package app

import (
	"context"

	"github.com/fishus/go-advanced-gophermart/internal/api"
	"github.com/fishus/go-advanced-gophermart/internal/service"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func RunAPIServer(ctx context.Context, storage store.Storager, loyalty AccrualDaemon) error {
	serviceConfig := &service.Config{
		JWTExpires:   Config.JWTExpires(),
		JWTSecretKey: Config.JWTSecretKey(),
	}
	serv := service.New(serviceConfig, storage)

	apiConfig := &api.Config{
		ServerAddr: Config.RunAddr(),
	}

	server := api.NewServer(apiConfig, serv, loyalty)
	Closers = append(Closers, server)
	return server.Run()
}
