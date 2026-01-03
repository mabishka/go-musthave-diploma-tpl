package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateShort(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		n    int
		want int
	}{
		{
			name: "length",
			n:    10,
			want: 10,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := CreateShort(test.n)

			assert.NoError(t, err)
			assert.Equal(t, test.want, len(got))

			got1, err := CreateShort(test.n)

			assert.NoError(t, err)
			assert.NotEqual(t, got, got1)
		})
	}
}
