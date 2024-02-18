// Package postgresclient ...
package postgresclient

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"trintech/review/config"
)

// PostgresClient ...
type PostgresClient struct {
	*pgxpool.Pool
	DB *config.Database
}

// Connect ...
func (c *PostgresClient) Connect(ctx context.Context) error {
	connectionString := c.DB.Address()
	u, err := url.Parse(connectionString)
	if err != nil {
		return fmt.Errorf("cannot create new connection to Postgres (failed to parse URI): %w", err)
	}
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return fmt.Errorf("cannot read PG_CONNECTION_URI: %w", err)
	}

	config.MaxConns = c.DB.MaxConnection
	config.MaxConnIdleTime = 15 * time.Second
	config.HealthCheckPeriod = 600 * time.Millisecond

	slog.Info("NewConnectionPool max connection", "max connection", c.DB.MaxConnection)

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("cannot create new connection to %q", u.Redacted()), err)
	}

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("unable to connect postgres: %w", err)
	}

	c.Pool = pool

	return nil
}

// Stop ...
func (c *PostgresClient) Stop(_ context.Context) error {
	c.Pool.Close()
	return nil
}
