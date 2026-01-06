package user

import (
	"context"
	"testing"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"github.com/mabishka/go-musthave-diploma-tpl/mock"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		conn    model.Connector
		want    *UserData
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

func TestUserData_AddUser(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		conn    model.Connector
		login   string
		token   string
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
			p := &UserData{conn: test.conn}

			err := p.AddUser(context.Background(), test.login, test.token)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestUserData_GetUser(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		conn    model.Connector
		token   string
		want    int
		want2   string
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
			p := &UserData{conn: test.conn}

			id, login, err := p.GetUser(context.Background(), test.token)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, id)
			assert.NotEmpty(t, login)
		})
	}
}

func TestUserData_GetAuth(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		conn    model.Connector
		login   string
		want    int
		want2   string
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
			p := &UserData{conn: test.conn}

			id, token, err := p.GetAuth(context.Background(), test.login)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, id)
			assert.NotEmpty(t, token)
		})
	}
}
