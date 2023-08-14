package config

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.uber.org/zap"
	"time"
)

// NewDatabaseConnection Initializes the connection pool for the database
func NewDatabaseConnection(log *zap.SugaredLogger, config DatabaseConfigurations) *bun.DB {
	sqlDB := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(config.DNS),
		pgdriver.WithTimeout(time.Duration(config.Timeout)*time.Second),
		pgdriver.WithDialTimeout(time.Duration(config.Timeout)*time.Second),
		pgdriver.WithReadTimeout(time.Duration(config.Timeout)*time.Second),
		pgdriver.WithWriteTimeout(time.Duration(config.Timeout)*time.Second),
	))

	sqlDB.SetMaxOpenConns(config.Pool)
	sqlDB.SetMaxIdleConns(config.Pool)
	sqlDB.SetConnMaxLifetime(time.Duration(config.Timeout) * time.Minute)

	// Test connect to DB
	err := sqlDB.Ping()
	if err != nil {
		log.Fatalf("Fail to connect to DB: %v", err)
	}

	db := bun.NewDB(sqlDB, pgdialect.New())

	log.Infof("Database connected. Connections opened: %d", db.Stats().OpenConnections)
	return db
}

// CloseDatabaseConnection Close the database connection
func CloseDatabaseConnection(database *bun.DB) {
	_ = database.Close()
}
