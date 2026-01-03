package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewToken(t *testing.T) {

	token, err := NewToken("user", "password")
	assert.NoError(t, err)
	tests := []struct {
		name string // description of this test case
		user     string
		password string
		want     string
		wantErr  bool
	}{
		{
			name:     "positive",
			user:     "user",
			password: "password",
			want:     token,
			wantErr:  false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewToken(test.user, test.password)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, token)
			assert.Equal(t, got, token)
		})
	}
}
