package handler

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
		addr    string
		want    *Server
		wantErr bool
	}{
		{
			name:    "positive",
			conn:    &mock.Connector{},
			addr:    "",
			wantErr: true,
		},
		{
			name:    "negative",
			conn:    &mock.ConnectorError{},
			addr:    "",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := New(context.Background(), test.conn, test.addr)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, got)
		})
	}
}
