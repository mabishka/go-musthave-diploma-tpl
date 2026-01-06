package db

import (
	"context"
	"testing"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"github.com/mabishka/go-musthave-diploma-tpl/mock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
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
			err := CreateUser(context.Background(), test.e)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestAddUser(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		e       model.Executor
		login   string
		auth    string
		wantErr bool
	}{
		{
			name:    "negative",
			e:       &mock.ExecutorError{},
			login:   "aaa",
			auth:    "aaa",
			wantErr: true,
		},
		{
			name:    "positive_empty",
			e:       &mock.Executor{},
			login:   "bbb",
			auth:    "bbb",
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := AddUser(context.Background(), test.e, test.login, test.auth)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestGetAuth(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		e       model.Executor
		login   string
		wantErr bool
	}{
		{
			name:    "negative",
			e:       &mock.ExecutorError{},
			login:   "aaa",
			wantErr: true,
		},
		{
			name:    "positive_empty",
			e:       &mock.Executor{},
			login:   "bbb",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			id, auth, err := GetAuth(context.Background(), test.e, test.login)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, id, 0)
			assert.Equal(t, auth, "")
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		e       model.Executor
		auth    string
		wantErr bool
	}{
		{
			name:    "negative",
			e:       &mock.ExecutorError{},
			auth:    "aaa",
			wantErr: true,
		},
		{
			name:    "positive_empty",
			e:       &mock.Executor{},
			auth:    "bbb",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			id, login, err := GetUser(context.Background(), test.e, test.auth)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, id, 0)
			assert.Equal(t, login, "")
		})
	}
}
