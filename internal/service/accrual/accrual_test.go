package accrual

import (
	"context"
	"net/http"
	"testing"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"github.com/mabishka/go-musthave-diploma-tpl/mock"
	"github.com/stretchr/testify/assert"
)

func TestAcuralData_request(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		conn model.Connector
		addr string
		data    model.AccrualProcess
		order   int
		want    model.OrderStateType
		wantErr bool
	}{
		{
			name:    "positive",
			addr:    "",
			conn:    &mock.Connector{},
			data:    model.AccrualProcess{Order: 1, User: 2},
			wantErr: true,
		},
		{
			name:    "negative",
			addr:    "",
			conn:    &mock.ConnectorError{},
			data:    model.AccrualProcess{Order: 1, User: 2},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := New(context.Background(), test.conn, test.addr)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, err := p.request(context.Background(), test.data.User, test.data.Order)
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, got, model.OrderStateFinish)
		})
	}
}

func TestAcuralData_StartProcess(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		conn    model.Connector
		addr    string
		data    model.AccrualProcess
		wantErr bool
	}{
		{
			name:    "positive",
			addr:    "",
			conn:    &mock.Connector{},
			data:    model.AccrualProcess{Order: 1, User: 2},
			wantErr: false,
		},
		{
			name:    "negative",
			addr:    "",
			conn:    &mock.ConnectorError{},
			data:    model.AccrualProcess{Order: 3, User: 4},
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, fn := context.WithCancel(context.Background())
			got, err := New(ctx, test.conn, test.addr)
			if test.wantErr {
				assert.Error(t, err)
				fn()
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, got)

			go got.StartProcess(ctx)
			got.ProcessOrder(test.data)
			fn()
		})
	}
}

func Test_parseResponceStatus(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		status  int
		want    model.OrderStateType
		wantErr bool
	}{
		{
			status:  http.StatusOK,
			want:    model.OrderStateFinish,
			wantErr: false,
		},
		{
			status:  http.StatusTooManyRequests,
			want:    model.OrderStateActive,
			wantErr: false,
		},
		{
			status:  http.StatusNoContent,
			want:    model.OrderStateActive,
			wantErr: false,
		},
		{
			status:  http.StatusInternalServerError,
			want:    model.OrderStateActive,
			wantErr: false,
		},
		{
			status:  http.StatusNotImplemented,
			want:    model.OrderStateActive,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseResponceStatus(test.status)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, got, test.want)
		})
	}
}

func Test_parseResponseDataStatus(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		status  model.AccrualStatusType
		want    model.OrderStateType
		wantErr bool
	}{
		{
			status:  model.AccrualStatusEmpty,
			want:    model.OrderStateFinish,
			wantErr: true,
		},
		{
			status:  model.AccrualStatusRegistered,
			want:    model.OrderStateActive,
			wantErr: false,
		},
		{
			status:  model.AccrualStatusInvalid,
			want:    model.OrderStateFinish,
			wantErr: true,
		},
		{
			status:  model.AccrualStatusProcessing,
			want:    model.OrderStateActive,
			wantErr: false,
		},
		{
			status:  model.AccrualStatusProcessed,
			want:    model.OrderStateFinish,
			wantErr: false,
		},
		{
			status:  "aaa",
			want:    model.OrderStateFinish,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseResponseDataStatus(test.status)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, got, test.want)
		})
	}
}
