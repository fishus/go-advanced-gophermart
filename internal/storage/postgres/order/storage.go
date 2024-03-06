package order

import (
	lStorage "github.com/fishus/go-advanced-gophermart/internal/storage/postgres/loyalty"
	"github.com/jackc/pgx/v5/pgxpool"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

type storage struct {
	cfg     *Config
	pool    *pgxpool.Pool
	loyalty store.Loyaltier
}

var _ store.Orderer = (*storage)(nil)

func New(pool *pgxpool.Pool, cfg *Config) (*storage, error) {
	loyalty, err := lStorage.New(pool, &lStorage.Config{
		ConnString:     cfg.ConnString,
		ConnectTimeout: cfg.ConnectTimeout,
		QueryTimeout:   cfg.QueryTimeout,
	})
	if err != nil {
		return nil, err
	}

	return &storage{
		cfg:     cfg,
		pool:    pool,
		loyalty: loyalty,
	}, nil
}

func (s *storage) Loyalty() store.Loyaltier {
	return s.loyalty
}
