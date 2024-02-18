package postgresclient

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"golang.org/x/net/context"
)

// Transaction ...
func (c *PostgresClient) Transaction(ctx context.Context, fn func(context.Context, pgx.Tx) error) error {
	tx, err := c.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return fmt.Errorf("unable to create transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := fn(ctx, tx); err != nil {
		return fmt.Errorf("unable to execute transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil
}
