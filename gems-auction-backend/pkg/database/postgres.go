package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresOptions struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	MaxConns int32
}

// NewPool creates a pgxpool.Pool with safe defaults.
func NewPool(opts PostgresOptions) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		opts.User,
		opts.Password,
		opts.Host,
		opts.Port,
		opts.DBName,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	if opts.MaxConns > 0 {
		cfg.MaxConns = opts.MaxConns
	}

	// pool health-check settings
	cfg.HealthCheckPeriod = 30 * time.Second
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.MaxConnLifetime = 60 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
