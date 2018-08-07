package main

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestNewCoinbaseTX(t *testing.T) {
	type args struct {
		to   string
		data string
	}
	tests := []struct {
		name string
		args args
		want *Transaction
	}{
		{
			name: "type check",
			args: args{
				to:   "To",
				data: "Data",
			},
			want: &Transaction{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCoinbaseTX(tt.args.to, tt.args.data)
			assert.IsType(t, tt.want, got, tt.name)
		})
	}
}
