package mock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnector_Executor(t *testing.T) {

	ctx := context.Background()
	query := "select 1"

	e := Executor{}

	_, err := e.QueryContext(ctx, query)
	assert.Error(t, err)

	result, err := e.ExecContext(ctx, query)
	assert.NoError(t, err)
	assert.Empty(t, result)

	err = e.Commit()
	assert.NoError(t, err)

	err = e.Rollback()
	assert.NoError(t, err)
}

func TestConnector_ExecutorError(t *testing.T) {

	ctx := context.Background()
	query := "select 1"

	e := ExecutorError{}

	_, err := e.QueryContext(ctx, query)
	assert.Error(t, err)

	result, err := e.ExecContext(ctx, query)
	assert.Error(t, err)
	assert.Empty(t, result)

	err = e.Commit()
	assert.NoError(t, err)

	err = e.Rollback()
	assert.NoError(t, err)
}

func TestConnector_Connector(t *testing.T) {

	ctx := context.Background()
	query := "select 1"

	e := Connector{}

	err := e.PingContext(ctx)
	assert.NoError(t, err)

	_, err = e.QueryContext(ctx, query)
	assert.Error(t, err)

	result, err := e.ExecContext(ctx, query)
	assert.NoError(t, err)
	assert.Empty(t, result)

	tx, err := e.BeginTx(ctx, nil)
	assert.NoError(t, err)
	assert.Empty(t, tx)
	assert.IsType(t, &Executor{}, tx)

	err = e.Commit()
	assert.NoError(t, err)

	err = e.Rollback()
	assert.NoError(t, err)

	err = e.Close()
	assert.NoError(t, err)
}

func TestConnector_ConnectorError(t *testing.T) {

	ctx := context.Background()
	query := "select 1"

	e := ConnectorError{}

	err := e.PingContext(ctx)
	assert.NoError(t, err)

	_, err = e.QueryContext(ctx, query)
	assert.Error(t, err)

	result, err := e.ExecContext(ctx, query)
	assert.Error(t, err)
	assert.Empty(t, result)

	tx, err := e.BeginTx(ctx, nil)
	assert.NoError(t, err)
	assert.Empty(t, tx)
	assert.IsType(t, &ExecutorError{}, tx)

	err = e.Commit()
	assert.NoError(t, err)

	err = e.Rollback()
	assert.NoError(t, err)

	err = e.Close()
	assert.NoError(t, err)
}
