package user

import (
	"github.com/jackc/pgx/v5/pgxpool"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

type storage struct {
	cfg  *Config
	pool *pgxpool.Pool
}

var _ store.Userer = (*storage)(nil)

func New(pool *pgxpool.Pool, cfg *Config) (*storage, error) {
	return &storage{
		cfg:  cfg,
		pool: pool,
	}, nil
}
