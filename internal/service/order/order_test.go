package order

import (
	"context"
	"testing"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"github.com/mabishka/go-musthave-diploma-tpl/mock"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		conn    model.Connector
		want    *OrderData
		wantErr bool
	}{
		{
			name:    "positive",
			conn:    &mock.Connector{},
			wantErr: false,
		},
		{
			name:    "negative",
			conn:    &mock.ConnectorError{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := New(context.Background(), test.conn)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, got)
		})
	}
}

func TestOrderData_Add(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		conn model.Connector
		order   int
		user    int
		wantErr bool
	}{
		{
			name:    "positive",
			conn:    &mock.Connector{},
			user:    1,
			order:   403017486,
			wantErr: false,
		},
		{
			name:    "negative luhn",
			conn:    &mock.Connector{},
			user:    2,
			order:   4,
			wantErr: true,
		},
		{
			name:    "negative db",
			conn:    &mock.ConnectorError{},
			user:    2,
			order:   403017486,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := OrderData{conn: test.conn}

			err := p.Add(context.Background(), test.order, test.user)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestOrderData_GetOrderList(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		conn model.Connector
		user    int
		want    []model.OrderResponse
		wantErr bool
	}{
		{
			name:    "positive",
			conn:    &mock.Connector{},
			user:    1,
			wantErr: true,
		},
		{
			name:    "negative",
			conn:    &mock.ConnectorError{},
			user:    2,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := OrderData{conn: test.conn}

			got, err := p.GetOrderList(context.Background(), test.user)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, len(got), 0)
		})
	}
}

func TestOrderData_GetActiveOrderList(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		conn    model.Connector
		want    []model.AccrualProcess
		wantErr bool
	}{
		{
			name:    "positive",
			conn:    &mock.Connector{},
			wantErr: true,
		},
		{
			name:    "negative",
			conn:    &mock.ConnectorError{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := OrderData{conn: test.conn}

			got, err := p.GetActiveOrderList(context.Background())
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, len(got), 0)
		})
	}
}

func TestOrderData_GetWithdrawals(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		conn model.Connector
		user    int
		want    []model.WithdrawnResponse
		wantErr bool
	}{
		{
			name:    "positive",
			conn:    &mock.Connector{},
			user:    1,
			wantErr: true,
		},
		{
			name:    "negative",
			conn:    &mock.ConnectorError{},
			user:    2,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := OrderData{conn: test.conn}

			got, err := p.GetWithdrawals(context.Background(), test.user)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, len(got), 0)
		})
	}
}

func TestOrderData_GetBalance(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		conn model.Connector
		user    int
		wantErr bool
	}{
		{
			name:    "positive",
			conn:    &mock.Connector{},
			user:    1,
			wantErr: true,
		},
		{
			name:    "negative",
			conn:    &mock.ConnectorError{},
			user:    2,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := OrderData{conn: test.conn}
			current, withdrawn, err := p.GetBalance(context.Background(), test.user)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			assert.InDelta(t, current, 0, 0.01)
			assert.InDelta(t, withdrawn, 0, 0.01)
		})
	}
}

func TestOrderData_Withdraw(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		conn model.Connector
		user    int
		order   int
		sum     float32
		wantErr bool
	}{
		{
			name:    "positive",
			conn:    &mock.Connector{},
			user:    1,
			order:   403017486,
			sum:     4,
			wantErr: true,
		},
		{
			name:    "negative luhn",
			conn:    &mock.Connector{},
			user:    2,
			order:   5,
			sum:     5.5,
			wantErr: true,
		},
		{
			name:    "negative db",
			conn:    &mock.ConnectorError{},
			user:    2,
			order:   403017486,
			sum:     5.5,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := OrderData{conn: test.conn}

			err := p.Withdraw(context.Background(), test.user, test.order, test.sum)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
