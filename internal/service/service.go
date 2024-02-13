package service

import (
	"context"

	oService "github.com/fishus/go-advanced-gophermart/internal/service/order"
	uService "github.com/fishus/go-advanced-gophermart/internal/service/user"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

type Userer interface {
	Register(context.Context, models.User) (models.UserID, error)
	Login(context.Context, models.User) (models.UserID, error)
	UserByID(context.Context, models.UserID) (*models.User, error)
	BuildToken(models.UserID) (string, error)
	DecryptToken(tokenString string) (*uService.JWTClaims, error)
	CheckAuthorizationHeader(auth string) (*uService.JWTClaims, error)
}

type Orderer interface {
	Add(ctx context.Context, userID models.UserID, orderNum string) (models.OrderID, error)
}

type service struct {
	cfg     *Config
	storage store.Storager
	user    Userer
	order   Orderer
}

func New(cfg *Config, s store.Storager) *service {
	order := oService.New(s)
	userCfg := &uService.Config{
		JWTExpires:   cfg.JWTExpires,
		JWTSecretKey: cfg.JWTSecretKey,
	}
	user := uService.New(userCfg, s)
	return &service{
		storage: s,
		cfg:     cfg,
		user:    user,
		order:   order,
	}
}

func (s *service) Storage() store.Storager {
	return s.storage
}

func (s *service) User() Userer {
	return s.user
}

func (s *service) Order() Orderer {
	return s.order
}
