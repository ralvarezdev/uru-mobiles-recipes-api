package sqlite

import (
	"log/slog"
	"time"

	godatabasessql "github.com/ralvarezdev/go-databases/sql"
	gojwtsyncsqlite "github.com/ralvarezdev/go-jwt/sync/sqlite"
)

const (
	// DriverName is the name of the SQLite driver
	DriverName = "sqlite3"

	// SyncDataSourceName is the data source name for JWT sync SQLite connection
	SyncDataSourceName = "file:sync.db?cache=shared&_journal_mode=WAL"

	// RabbitMQConsumerDataSourceName is the data source name for RabbitMQ SQLite connection
	RabbitMQConsumerDataSourceName = "file:rabbitmq_consumer.db?cache=shared&_journal_mode=WAL"

	// MaxOpenConnections is the maximum number of open connections to the SQLite database
	MaxOpenConnections = 10

	// MaxIdleConnections is the maximum number of idle connections to the SQLite database
	MaxIdleConnections = 5

	// ConnectionMaxLifetime is the maximum amount of time a connection may be reused
	ConnectionMaxLifetime = 30 * time.Minute

	// ConnectionMaxIdleTime is the maximum amount of time a connection may be idle
	ConnectionMaxIdleTime = 10 * time.Minute
)

var (
	// SyncConfig is the JWT sync config
	SyncConfig = godatabasessql.Config{
		DriverName:            DriverName,
		DataSourceName:        SyncDataSourceName,
		MaxOpenConnections:    MaxOpenConnections,
		MaxIdleConnections:    MaxIdleConnections,
		ConnectionMaxLifetime: ConnectionMaxLifetime,
		ConnectionMaxIdleTime: ConnectionMaxIdleTime,
	}

	// TokenValidatorConfig is the token validator config
	TokenValidatorConfig = godatabasessql.Config{
		DriverName:            DriverName,
		DataSourceName:        RabbitMQConsumerDataSourceName,
		MaxOpenConnections:    MaxOpenConnections,
		MaxIdleConnections:    MaxIdleConnections,
		ConnectionMaxLifetime: ConnectionMaxLifetime,
		ConnectionMaxIdleTime: ConnectionMaxIdleTime,
	}

	// SyncSQLiteService is the JWT sync SQLite service
	SyncSQLiteService godatabasessql.Service

	// SyncService is the JWT sync service
	SyncService *gojwtsyncsqlite.Service

	// TokenValidatorService is the JWT token validator SQLite service
	TokenValidatorService godatabasessql.Service
)

// Load initializes the SQLite handlers and services
//
// Parameters:
//
//   - logger: The logger (optional, can be nil)
func Load(logger *slog.Logger) {
	// Initialize the JWT sync SQLite service
	syncSQLiteService, err := godatabasessql.NewDefaultService(
		&SyncConfig,
	)
	if err != nil {
		panic(err)
	}
	SyncSQLiteService = syncSQLiteService

	// Initialize the JWT sync service
	syncService, err := gojwtsyncsqlite.NewService(
		SyncSQLiteService,
		logger,
	)
	if err != nil {
		panic(err)
	}
	SyncService = syncService

	// Initialize the token validator SQLite service
	tokenValidatorService, err := godatabasessql.NewDefaultService(
		&TokenValidatorConfig,
	)
	if err != nil {
		panic(err)
	}
	TokenValidatorService = tokenValidatorService
}
