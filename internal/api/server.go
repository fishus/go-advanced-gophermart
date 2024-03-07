package api

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	apiLoyalty "github.com/fishus/go-advanced-gophermart/internal/api/loyalty"
	apiOrder "github.com/fishus/go-advanced-gophermart/internal/api/order"
	apiUser "github.com/fishus/go-advanced-gophermart/internal/api/user"
	"github.com/fishus/go-advanced-gophermart/internal/logger"
	"github.com/fishus/go-advanced-gophermart/internal/service"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=Servicer --with-expecter
type Servicer interface {
	User() service.Userer
	Order() service.Orderer
	Loyalty() service.Loyaltier
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=AccrualDaemon  --with-expecter
type AccrualDaemon interface {
	AddNewOrder(context.Context, models.Order)
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=User --with-expecter
type User interface {
	Login(http.ResponseWriter, *http.Request)
	Register(http.ResponseWriter, *http.Request)
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=Order --with-expecter
type Order interface {
	Add(http.ResponseWriter, *http.Request)
	List(http.ResponseWriter, *http.Request)
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=Loyalty --with-expecter
type Loyalty interface {
	Balance(http.ResponseWriter, *http.Request)
	Withdraw(http.ResponseWriter, *http.Request)
	Withdrawals(http.ResponseWriter, *http.Request)
}

type server struct {
	cfg     *Config
	server  *http.Server
	service Servicer
	daemon  AccrualDaemon
	user    User
	order   Order
	loyalty Loyalty
}

func NewServer(cfg *Config, service Servicer, daemon AccrualDaemon) *server {
	s := &server{
		cfg:     cfg,
		service: service,
		daemon:  daemon,
		user:    apiUser.NewAPI(service),
		order:   apiOrder.NewAPI(service, daemon),
		loyalty: apiLoyalty.NewAPI(service),
	}

	srv := &http.Server{
		Addr:              cfg.ServerAddr,
		Handler:           s.Router(),
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
