package luhn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_checksum(t *testing.T) {
	tests := []struct {
		name   string // description of this test case
		number int
		want   int
	}{
		{
			number: 123,
			want:   0,
		},
		{
			number: 403017486,
			want:   2,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := checksum(test.number)
			assert.Equal(t, got, test.want)
		})
	}
}

func TestValid(t *testing.T) {
	tests := []struct {
		name   string // description of this test case
		number int
		want   bool
	}{
		{
			number: 123,
			want:   false,
		},
		{
			number: 403017486,
			want:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Valid(test.number)
			assert.Equal(t, got, test.want)
		})
	}
}
