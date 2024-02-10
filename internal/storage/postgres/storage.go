package postgres

import (
	"context"
	"io"

	"github.com/jackc/pgx/v5/pgxpool"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

type storage struct {
	cfg  *Config
	pool *pgxpool.Pool
}

func (s *storage) Close() error {
	s.pool.Close()
	return nil
}

var _ io.Closer = (*storage)(nil)
var _ store.Storager = (*storage)(nil)

func New(ctx context.Context, cfg *Config) (store.Storager, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()
	pool, err := pgxpool.New(ctx, cfg.ConnString)
	if err != nil {
		return nil, err
	}

	return &storage{cfg: cfg, pool: pool}, nil
}
