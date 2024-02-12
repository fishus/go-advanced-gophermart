package app

import (
	"context"
	"github.com/fishus/go-advanced-gophermart/internal/api"
	"github.com/fishus/go-advanced-gophermart/internal/service"
)

func RunAPIServer(ctx context.Context) error {
	db, err := ConnDB(ctx)
	if err != nil {
		return err
	}

	serviceConfig := &service.Config{
		JWTExpires:   Config.jwtExpires,
		JWTSecretKey: Config.jwtSecretKey,
	}
	serv := service.New(serviceConfig, db)

	apiConfig := &api.Config{
		ServerAddr: Config.RunAddr(),
	}

	server := api.NewServer(apiConfig, serv)
	Closers = append(Closers, server)
	err = server.Run()
	return err
}
