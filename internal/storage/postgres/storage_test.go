package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresTestSuite struct {
	suite.Suite
	storage

	tc  *tcpostgres.PostgresContainer
	cfg *Config
}

func (ts *PostgresTestSuite) SetupSuite() {
	const (
		dbUsername = "postgres"
		dbPassword = "password"
		dbName     = "postgres"
	)

	var (
		err    error
		dbHost string
		dbPort uint16
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pgc, err := tcpostgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16"),
		tcpostgres.WithDatabase(dbName),
		tcpostgres.WithUsername(dbUsername),
		tcpostgres.WithPassword(dbPassword),
		tcpostgres.WithInitScripts(),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)

	ts.Require().NoError(err)

	dbHost, err = pgc.Host(ctx)
	ts.Require().NoError(err)

	port, err := pgc.MappedPort(ctx, "5432")
	ts.Require().NoError(err)

	dbPort = uint16(port.Int())

	cfg := &Config{
		ConnectTimeout: 5 * time.Second,
		QueryTimeout:   5 * time.Second,
		ConnString:     fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbUsername, dbPassword, dbHost, dbPort, dbName),
	}

	ts.tc = pgc
	ts.cfg = cfg

	db, err := New(ctx, cfg)
	ts.Require().NoError(err)

	ts.storage = *db

	ts.T().Logf("Stared postgres at %s:%d", dbHost, dbPort)
}

func (ts *PostgresTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ts.Require().NoError(ts.tc.Terminate(ctx))
}

func (ts *PostgresTestSuite) SetupTest() {
	ts.Require().NoError(ts.clean(context.Background()))
}

func (ts *PostgresTestSuite) TearDownTest() {
	ts.Require().NoError(ts.clean(context.Background()))
}

func TestPostgres(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}

func (s *storage) clean(ctx context.Context) (err error) {
	newCtx, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	_, err = s.pool.Exec(newCtx, "TRUNCATE TABLE users CASCADE;")
	if err != nil {
		return
	}

	_, err = s.pool.Exec(newCtx, "TRUNCATE TABLE orders CASCADE;")
	if err != nil {
		return
	}

	_, err = s.pool.Exec(newCtx, "TRUNCATE TABLE loyalty_history CASCADE;")
	if err != nil {
		return
	}

	_, err = s.pool.Exec(newCtx, "TRUNCATE TABLE loyalty_balance CASCADE;")
	if err != nil {
		return
	}

	return
}
