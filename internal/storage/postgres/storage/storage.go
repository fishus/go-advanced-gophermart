package storage

import (
	"context"
	"io"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	lStorage "github.com/fishus/go-advanced-gophermart/internal/storage/postgres/loyalty"
	"github.com/fishus/go-advanced-gophermart/internal/storage/postgres/migration"
	oStorage "github.com/fishus/go-advanced-gophermart/internal/storage/postgres/order"
	uStorage "github.com/fishus/go-advanced-gophermart/internal/storage/postgres/user"
)

type storage struct {
	cfg     *Config
	pool    *pgxpool.Pool
	user    store.Userer
	order   store.Orderer
	loyalty store.Loyaltier
}

func (s *storage) Close() error {
	logger.Log.Info("Shutdown DB pool")
	s.pool.Close()
	return nil
}

var _ io.Closer = (*storage)(nil)
var _ store.Storager = (*storage)(nil)

func New(ctx context.Context, cfg *Config) (*storage, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	pgxConfig, err := pgxpool.ParseConfig(cfg.ConnString)
	if err != nil {
		return nil, err
	}

	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		// Add integration with decimal package
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, err
	}

	if err := migration.Migrate(pool); err != nil {
		return nil, err
	}

	user, err := uStorage.New(pool, &uStorage.Config{
		ConnString:     cfg.ConnString,
		ConnectTimeout: cfg.ConnectTimeout,
		QueryTimeout:   cfg.QueryTimeout,
	})
	if err != nil {
		return nil, err
	}

	order, err := oStorage.New(pool, &oStorage.Config{
		ConnString:     cfg.ConnString,
		ConnectTimeout: cfg.ConnectTimeout,
		QueryTimeout:   cfg.QueryTimeout,
	})
	if err != nil {
		return nil, err
	}

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
		user:    user,
		order:   order,
		loyalty: loyalty,
	}, nil
}

func (s *storage) User() store.Userer {
	return s.user
}

func (s *storage) Order() store.Orderer {
	return s.order
}

func (s *storage) Loyalty() store.Loyaltier {
	return s.loyalty
}
