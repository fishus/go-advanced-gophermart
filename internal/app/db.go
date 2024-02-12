package app

import (
	"context"
	"time"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	"github.com/fishus/go-advanced-gophermart/internal/storage/postgres"
)

func ConnDB(ctx context.Context) (store.Storager, error) {
	dbConfig := &postgres.Config{
		ConnString:     Config.DatabaseURI(),
		ConnectTimeout: 5 * time.Second,
		QueryTimeout:   5 * time.Second,
	}

	db, err := postgres.New(ctx, dbConfig)
	if err != nil {
		return nil, err
	}
	Closers = append(Closers, db)
	return db, nil
}
