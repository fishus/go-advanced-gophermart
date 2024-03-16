package loyalty

import (
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

type Orderer interface {
	ValidateNumLuhn(num string) error
}

type service struct {
	storage store.Storager
	order   Orderer
}

func New(s store.Storager) *service {
	return &service{storage: s}
}

func (s *service) Storage() store.Storager {
	return s.storage
}

func (s *service) SetOrder(order Orderer) *service {
	s.order = order
	return s
}
