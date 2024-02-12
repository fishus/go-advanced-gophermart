package service

import (
	"context"

	uService "github.com/fishus/go-advanced-gophermart/internal/service/user"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

type Userer interface {
	Register(context.Context, models.User) (models.UserID, error)
	Login(context.Context, models.User) (models.UserID, error)
	UserByID(context.Context, models.UserID) (*models.User, error)
	BuildToken(user models.UserID) (string, error)
}

type service struct {
	cfg     *Config
	storage store.Storager
	user    Userer
}

func New(cfg *Config, s store.Storager) *service {
	userCfg := &uService.Config{
		JWTExpires:   cfg.JWTExpires,
		JWTSecretKey: cfg.JWTSecretKey,
	}
	user := uService.New(userCfg, s)
	return &service{
		storage: s,
		cfg:     cfg,
		user:    user,
	}
}

func (s *service) Storage() store.Storager {
	return s.storage
}

func (s *service) User() Userer {
	return s.user
}
