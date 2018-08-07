package main

import (
	"os"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestNewBlockchain(t *testing.T) {
	tests := []struct {
		name    string
		want    *Blockchain
		wantErr bool
	}{
		{
			name: "create a chain",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBlockchain("input")
			assert.NoError(t, err, "no error")
			assert.IsType(t, &Blockchain{}, got, "if there's no error, this should be of type *Blockchain")
		})
	}
	os.Remove(dbFile)
}

func TestBlockchain_AddBlock(t *testing.T) {
	targetBits = 8
	os.Remove(dbFile)
	bc, _ := CreateBlockchain("address")

	err := bc.AddBlock([]*Transaction{&Transaction{}})
	assert.NoError(t, err, "this should persist a transaction")
	err = bc.AddBlock([]*Transaction{&Transaction{}})
	assert.NoError(t, err, "this should persist another transaction")
	os.Remove(dbFile)
}

func BenchmarkBlockchain_AddBlock(b *testing.B) {
	bc, _ := CreateBlockchain("address")
	for n := 0; n < b.N; n++ {
		bc.AddBlock([]*Transaction{})
	}
	os.Remove(dbFile)
}
