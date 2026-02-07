package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sfdeloach/churchsite/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB wraps the PostgreSQL and Redis connections.
type DB struct {
	Postgres *gorm.DB
	Redis    *redis.Client
}

// Connect establishes connections to PostgreSQL and Redis.
func Connect(cfg *config.Config) (*DB, error) {
	logLevel := logger.Warn
	if cfg.IsDevelopment() {
		logLevel = logger.Info
	}

	pg, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := pg.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	opts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(opts)

	db := &DB{
		Postgres: pg,
		Redis:    rdb,
	}

	if err := db.PingPostgres(); err != nil {
		return nil, err
	}
	slog.Info("connected to PostgreSQL")

	if err := db.PingRedis(); err != nil {
		return nil, err
	}
	slog.Info("connected to Redis")

	return db, nil
}

// PingPostgres checks the PostgreSQL connection.
func (db *DB) PingPostgres() error {
	sqlDB, err := db.Postgres.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// PingRedis checks the Redis connection.
func (db *DB) PingRedis() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return db.Redis.Ping(ctx).Err()
}

// Close shuts down both database connections.
func (db *DB) Close() {
	if sqlDB, err := db.Postgres.DB(); err == nil {
		sqlDB.Close()
	}
	db.Redis.Close()
}
