package main

import (
	"context"
	"testing"
	"time"
)

func Test_run(t *testing.T) {

	ctx, fnCancel := context.WithTimeout(context.Background(), time.Second*2)
	defer fnCancel()
	tests := []struct {
		name    string // description of this test case
		wantErr bool
	}{
		{
			name:    "positive",
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			run(ctx)
		})
	}

}
