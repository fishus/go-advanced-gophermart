package user

import (
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

type service struct {
	cfg     *Config
	storage store.Storager
}

func New(cfg *Config, s store.Storager) *service {
	return &service{cfg: cfg, storage: s}
}

func (s *service) Storage() store.Storager {
	return s.storage
}
