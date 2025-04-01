package pgxutil

import (
	"context"
	"fmt"
	"sync"

	"github.com/wagecloud/wagecloud-server/config"
	"github.com/wagecloud/wagecloud-server/internal/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var pgxPool *pgxpool.Pool
var m sync.Mutex

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)

	Begin(context.Context) (pgx.Tx, error)
}

// InitPgConnectionPool the database connection pgxPool.
func InitPgConnectionPool(postgresConfig config.Postgres) error {
	m.Lock()
	defer m.Unlock()

	if pgxPool != nil {
		return nil // The connection pgxPool has already been initialized
	}

	connStr := GetConnStr(postgresConfig)

	connConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		fmt.Println("Failed to parse config:", err)
		return err
	}

	// Set maximum number of connections
	connConfig.MaxConns = postgresConfig.MaxConnections
	// connConfig.MaxConnIdleTime = time.Duration(postgresConfig.MaxConnIdleTime) * time.Minute
	connConfig.ConnConfig.OnNotice = func(conn *pgconn.PgConn, notice *pgconn.Notice) {
		logger.Log.Info("notice", zap.String("message", notice.Message))
	}

	pgxPool, err = pgxpool.NewWithConfig(context.Background(), connConfig)

	if err != nil {
		return err
	}
	return nil
}

func GetPgxPool() (*pgxpool.Pool, error) {
	if pgxPool != nil {
		return pgxPool, nil
	}

	err := InitPgConnectionPool(config.GetConfig().Postgres)

	if err != nil {
		return nil, err
	}

	return pgxPool, nil
}

func GetPgxConn(ctx context.Context) (*pgxpool.Conn, error) {
	pgxPool, err := GetPgxPool()
	if err != nil {
		return nil, err
	}

	conn, err := pgxPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func InitSchema(ctx context.Context, postgresConfig config.Postgres, schema string) (err error) {
	connStr := GetConnStr(postgresConfig)

	pgConn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return err
	}
	defer pgConn.Close(ctx)

	// Create schema if it doesn't exist
	// Ignore error if schema already exists or if the user doesn't have permission to create schema
	pgConn.Exec(
		ctx,
		fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, schema),
	)

	// Set search path to schema so that we don't have to specify the schema name
	_, err = pgConn.Exec(
		ctx,
		fmt.Sprintf(`SET search_path TO %s`, schema),
	)
	if err != nil {
		return err
	}

	return nil
}

// Close the database connection pgxPool.
func ClosePgxPool() {
	m.Lock()
	defer m.Unlock()

	if pgxPool != nil {
		pgxPool.Close()
	}
}

func GetConnStr(postgresConfig config.Postgres) string {
	if postgresConfig.Url == "" {
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			postgresConfig.Host,
			postgresConfig.Port,
			postgresConfig.Username,
			postgresConfig.Password,
			postgresConfig.Database,
		)
	}

	return postgresConfig.Url
}
