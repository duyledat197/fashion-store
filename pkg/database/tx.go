package database

import (
	"context"
	"database/sql"
)

// Transaction implements a passing function with parameter have pointer of sql.Tx.
// The transaction begin with serializable isolation and then call passing function and then commit or rollback.
func Transaction(ctx context.Context, db Database, fn func(ctx context.Context, db *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := fn(ctx, tx); err != nil {
		return err
	}

	tx.Commit()

	return nil
}
