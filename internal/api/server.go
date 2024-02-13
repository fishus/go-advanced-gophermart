package api

import (
	"context"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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

func Router(s *server) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(9, "application/json"))
	r.Use(middleware.RequestLogger(&logger.LogFormatter{}))

	r.Post("/api/user/register", s.userRegister) // Регистрация пользователя
	r.Post("/api/user/login", s.userLogin)       // Аутентификация пользователя
	r.Post("/api/user/orders", s.orderAdd)       // Загрузка номера заказа для расчёта

	return r
}

func (s *server) Run() error {
	logger.Log.Info("Running api server", logger.String("address", s.cfg.ServerAddr))
	return s.server.ListenAndServe()
}

func (s *server) Close() error {
	logger.Log.Info("Shutdown api server")
	ctx := context.Background()
	return s.server.Shutdown(ctx)
}

var _ io.Closer = (*server)(nil)
