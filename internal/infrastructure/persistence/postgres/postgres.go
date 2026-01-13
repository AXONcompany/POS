package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/AXONcompany/POS/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func Connect(ctx context.Context, cfg config.Config) (*DB, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.GetPostgresDSN())
	if err != nil {
		return nil, fmt.Errorf("parse postgres dsn: %w", err)
	}

	poolCfg.MaxConns = cfg.PGMaxConns

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &DB{Pool: pool}, nil

}

func (d *DB) Close() {
	if d != nil && d.Pool != nil {
		d.Pool.Close()
	}
}
