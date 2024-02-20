package database

import (
	"context"
	"database/sql"

	"trintech/review/pkg/processor"
)

// Executor is a presentation of an database executor with exec and query command.
// Example implementation is [database/sql.Tx] and [database/sql.DB]
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Database is presentation of an database ase driver which implement [Executor], [trintech/review/pkg/processor.Processor] and transaction.
type Database interface {
	Executor
	processor.Factory
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}
