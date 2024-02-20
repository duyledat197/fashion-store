package postgres_client

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// PostgresClient is presentation for a custom client of postgres with [database/sql] based.
type PostgresClient struct {
	*sql.DB
	connectionString string
}

// New creates a new PostgresClient using the given connection string.
func NewPostgresClient(connString string) *PostgresClient {
	return &PostgresClient{
		connectionString: connString,
	}
}

// Connect implements postgres connection by [PostgresClient].
func (c *PostgresClient) Connect(ctx context.Context) error {
	var err error
	c.DB, err = sql.Open("postgres", c.connectionString)
	if err != nil {
		return err
	}

	if err := c.DB.Ping(); err != nil {
		return err
	}

	log.Println("connect postgres successful")

	return nil
}

// Close implements close postgres connection by [PostgresClient]..
func (c *PostgresClient) Close(ctx context.Context) error {
	return c.DB.Close()
}
