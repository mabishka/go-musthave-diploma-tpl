package model

import (
	"context"
	"database/sql"
)

type Creator interface {
	Executor
	PingContext(ctx context.Context) error
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Executor, error)
}

type Connector interface {
	Executor
	PingContext(ctx context.Context) error
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Executor, error)
	Close() error
}

type Executor interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	Commit() error
	Rollback() error
}

type DB struct {
	*sql.DB
}

func (p *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (Executor, error) {
	return p.DB.BeginTx(ctx, opts)
}
func (p *DB) Commit() error {
	return nil
}
func (p *DB) Rollback() error {
	return nil
}
