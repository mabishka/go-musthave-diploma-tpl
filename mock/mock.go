package mock

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
)

type Connector struct {
}

func (p *Connector) QueryContext(context.Context, string, ...any) (*sql.Rows, error) {
	return nil, errors.New("unsupported")
}

func (p *Connector) ExecContext(context.Context, string, ...any) (sql.Result, error) {
	return nil, nil
}

func (p *Connector) PingContext(ctx context.Context) error {
	return nil
}
func (p *Connector) BeginTx(ctx context.Context, opts *sql.TxOptions) (model.Executor, error) {
	return &Executor{}, nil
}
func (p *Connector) Close() error {
	return nil
}
func (p *Connector) Commit() error {
	return nil
}
func (p *Connector) Rollback() error {
	return nil
}

type ConnectorError struct {
}

func (p *ConnectorError) QueryContext(context.Context, string, ...any) (*sql.Rows, error) {
	return nil, errors.New("unsupported")
}

func (p *ConnectorError) ExecContext(context.Context, string, ...any) (sql.Result, error) {
	return nil, errors.New("unsupported")
}

func (p *ConnectorError) PingContext(ctx context.Context) error {
	return nil
}
func (p *ConnectorError) BeginTx(ctx context.Context, opts *sql.TxOptions) (model.Executor, error) {
	return &ExecutorError{}, nil
}
func (p *ConnectorError) Close() error {
	return nil
}
func (p *ConnectorError) Commit() error {
	return nil
}
func (p *ConnectorError) Rollback() error {
	return nil
}

type Executor struct {
}

func (p *Executor) QueryContext(context.Context, string, ...any) (*sql.Rows, error) {
	return nil, errors.New("unsupported")
}

func (p *Executor) ExecContext(context.Context, string, ...any) (sql.Result, error) {
	return nil, nil
}

func (p *Executor) Commit() error {
	return nil
}
func (p *Executor) Rollback() error {
	return nil
}

type ExecutorError struct {
}

func (p *ExecutorError) QueryContext(context.Context, string, ...any) (*sql.Rows, error) {
	return nil, errors.New("unsupported")
}

func (p *ExecutorError) ExecContext(context.Context, string, ...any) (sql.Result, error) {
	return nil, errors.New("unsupported")
}

func (p *ExecutorError) Commit() error {
	return nil
}
func (p *ExecutorError) Rollback() error {
	return nil
}
