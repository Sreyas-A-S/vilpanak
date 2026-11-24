package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"vilpanak/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect(cfg *config.Config) (*pgxpool.Pool, error) {
	connStr := cfg.DBConnectionString()
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	// Performance tuning pool settings
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Printf("Warning: Failed to ping PostgreSQL database (%v). App will run, but DB queries will fail until DB is online.", err)
	} else {
		log.Println("Connected to PostgreSQL successfully.")
	}

	Pool = pool
	return pool, nil
}
