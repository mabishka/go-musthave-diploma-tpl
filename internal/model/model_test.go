package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOrderStatus(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		x    AccrualStatusType
		want OrderStatusType
	}{
		{
			x:    AccrualStatusEmpty,
			want: OrderStatusNew,
		},
		{
			x:    AccrualStatusRegistered,
			want: OrderStatusNew,
		},
		{
			x:    AccrualStatusInvalid,
			want: OrderStatusInvalid,
		},
		{
			x:    AccrualStatusProcessing,
			want: OrderStatusProcessing,
		},
		{
			x:    AccrualStatusProcessed,
			want: OrderStatusProcessed,
		},
		{
			x:    "aa",
			want: OrderStatusUndefined,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GetOrderStatus(test.x)
			assert.Equal(t, got, test.want)
		})
	}
}

func TestGetState(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		x    AccrualStatusType
		want OrderStateType
	}{
		{
			x:    AccrualStatusEmpty,
			want: OrderStateUndefined,
		},
		{
			x:    AccrualStatusRegistered,
			want: OrderStateActive,
		},
		{
			x:    AccrualStatusInvalid,
			want: OrderStateFinish,
		},
		{
			x:    AccrualStatusProcessing,
			want: OrderStateActive,
		},
		{
			x:    AccrualStatusProcessed,
			want: OrderStateFinish,
		},
		{
			x:    "aa",
			want: OrderStateUndefined,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GetState(test.x)
			assert.Equal(t, got, test.want)
		})
	}
}
