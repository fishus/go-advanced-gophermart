package user

import (
	"context"
	"fmt"
	"github.com/fishus/go-advanced-gophermart/internal/storage/postgres/migration"
	"testing"
	"time"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
		dbName     = "users"
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

	pool, err := ts.conn(ctx)
	ts.Require().NoError(err)

	db, err := New(pool, cfg)
	ts.Require().NoError(err)

	ts.storage = *db

	err = migration.Migrate(pool)
	ts.Require().NoError(err)
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

func (ts *PostgresTestSuite) conn(ctx context.Context) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, ts.cfg.ConnectTimeout)
	defer cancel()

	pgxConfig, err := pgxpool.ParseConfig(ts.cfg.ConnString)
	if err != nil {
		return nil, err
	}

	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		// Add integration with decimal package
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}

	return pgxpool.NewWithConfig(ctx, pgxConfig)
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

//func (s *storage) addTestUser(ctx context.Context) (userID models.UserID, err error) {
//	bUsername := make([]byte, 10)
//	_, err = rand.Read(bUsername)
//	if err != nil {
//		return
//	}
//
//	userData := models.User{
//		Username: hex.EncodeToString(bUsername),
//		Password: hex.EncodeToString(bUsername),
//	}
//	return s.UserAdd(ctx, userData)
//}
