package migration

import (
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations
var migrations embed.FS

func Migrate(pool *pgxpool.Pool) (err error) {
	goose.SetBaseFS(migrations)

	if err = goose.SetDialect("postgres"); err != nil {
		return
	}

	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	err = goose.Up(db, "migrations")
	return
}
