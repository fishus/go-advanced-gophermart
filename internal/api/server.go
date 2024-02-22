package api

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
	"github.com/fishus/go-advanced-gophermart/internal/service"
)

type Servicer interface {
	User() service.Userer
	Order() service.Orderer
}

type server struct {
	cfg     *Config
	server  *http.Server
	service Servicer
}

func NewServer(cfg *Config, service Servicer) *server {
	s := &server{
		cfg:     cfg,
		service: service,
	}

	srv := &http.Server{
		Addr:              cfg.ServerAddr,
		Handler:           Router(s),
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}
	s.server = srv

	return s
}

func (s *server) Run() error {
	logger.Log.Info("Running api server", logger.String("address", s.cfg.ServerAddr))
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *server) Close() error {
	logger.Log.Info("Shutdown api server")
	ctx, cancel := context.WithTimeout(context.Background(), (15 * time.Second))
	go func() {
		<-ctx.Done()
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			logger.Log.Error("shutdown api server timed out.. forcing exit.")
		}
	}()
	err := s.server.Shutdown(ctx)
	cancel()
	return err
}

var _ io.Closer = (*server)(nil)
