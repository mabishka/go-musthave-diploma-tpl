package db

import (
	"context"
	"testing"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		addr    string
		want    model.Connector
		wantErr bool
	}{
		{
			name:    "negative",
			addr:    "",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conn, err := New(context.Background(), test.addr)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, conn, nil)
		})
	}
}
