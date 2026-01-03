package db

import (
	"context"
	"testing"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"github.com/mabishka/go-musthave-diploma-tpl/mock"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		e       model.Executor
		wantErr bool
	}{
		{
			name:    "negative",
			e:       &mock.ExecutorError{},
			wantErr: true,
		},
		{
			name:    "positive_empty",
			e:       &mock.Executor{},
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CreateOrder(context.Background(), test.e)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestProcessOrder(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		e       model.Executor
		order   int
		user    int
		wantErr bool
	}{
		{
			name:    "negative",
			e:       &mock.ExecutorError{},
			order:   1,
			user:    1,
			wantErr: true,
		},
		{
			name:    "positive_empty",
			e:       &mock.Executor{},
			order:   2,
			user:    2,
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ProcessOrder(context.Background(), test.e, test.order, test.user)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestGetBalance(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		e       model.Executor
		user    int
		wantErr bool
	}{
		{
			name:    "negative",
			e:       &mock.ExecutorError{},
			user:    1,
			wantErr: true,
		},
		{
			name:    "positive_empty",
			e:       &mock.Executor{},
			user:    2,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			current, withdrawn, err := GetBalance(context.Background(), test.e, test.user)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, current, 0)
			assert.Equal(t, withdrawn, 0)
		})
	}
}

func TestSetStatus(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		e       model.Executor
		order   int
		status  model.AccrualStatusType
		wantErr bool
	}{
		{
			name:    "negative",
			e:       &mock.ExecutorError{},
			order:   1,
			status:  model.AccrualStatusProcessed,
			wantErr: true,
		},
		{
			name:    "positive_empty",
			e:       &mock.Executor{},
			order:   2,
			status:  model.AccrualStatusInvalid,
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := SetStatus(context.Background(), test.e, test.order, test.status)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestAccrual(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		e       model.Executor
		user    int
		order   int
		value   float32
		wantErr bool
	}{
		{
			name:    "negative",
			e:       &mock.ExecutorError{},
			order:   1,
			value:   5,
			wantErr: true,
		},
		{
			name:    "positive_empty",
			e:       &mock.Executor{},
			order:   2,
			value:   5.5,
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Accrual(context.Background(), test.e, test.user, test.order, test.value)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestWithdraw(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		e       model.Executor
		user    int
		order   int
		value   float32
		wantErr bool
	}{
		{
			name:    "negative",
			e:       &mock.ExecutorError{},
			order:   1,
			value:   5,
			wantErr: true,
		},
		{
			name:    "positive_empty",
			e:       &mock.Executor{},
			order:   2,
			value:   5.5,
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Withdraw(context.Background(), test.e, test.user, test.order, test.value)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestGetWithdrawals(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		e       model.Executor
		user    int
		want    []model.WithdrawnResponse
		wantErr bool
	}{
		{
			name:    "negative",
			e:       &mock.ExecutorError{},
			user:    1,
			wantErr: true,
		},
		{
			name:    "positive_empty",
			e:       &mock.Executor{},
			user:    2,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			list, err := GetWithdrawals(context.Background(), test.e, test.user)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, len(list), 0)
		})
	}
}

func TestGetOrderList(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		e       model.Executor
		user    int
		want    []model.OrderResponse
		wantErr bool
	}{
		{
			name:    "negative",
			e:       &mock.ExecutorError{},
			user:    1,
			wantErr: true,
		},
		{
			name:    "positive_empty",
			e:       &mock.Executor{},
			user:    2,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			list, err := GetOrderList(context.Background(), test.e, test.user)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, len(list), 0)
		})
	}
}

func TestGetActiveOrderList(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		e       model.Executor
		want    []model.AccrualProcess
		wantErr bool
	}{
		{
			name:    "negative",
			e:       &mock.ExecutorError{},
			wantErr: true,
		},
		{
			name:    "positive_empty",
			e:       &mock.Executor{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			list, err := GetActiveOrderList(context.Background(), test.e)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, len(list), 0)
		})
	}
}
