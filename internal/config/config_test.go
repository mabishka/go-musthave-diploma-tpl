package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		{
			name: "positive",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := New()
			assert.Equal(t, got, New())
		})
	}
}

func TestConfig_Load(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		{
			name: "positive",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := New()
			c.Load()
		})
	}
}
