package user

import (
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

type Orderer interface {
	ValidateNumLuhn(num string) error
}

type service struct {
	cfg     *Config
	storage store.Storager
	order   Orderer
}

func New(cfg *Config, s store.Storager) *service {
	return &service{cfg: cfg, storage: s}
}

func (s *service) Storage() store.Storager {
	return s.storage
}

func (s *service) SetOrder(order Orderer) *service {
	s.order = order
	return s
}
