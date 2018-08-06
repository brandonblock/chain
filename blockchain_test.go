package chain

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

// func TestBlockchain_AddBlock(t *testing.T) {
// 	bc, _ := NewBlockchain()
// assert.NoError(t, err, "blockchain creation should succeed")

// bc.AddBlock("send 1 coin to someone")
// bc.AddBlock("send 1 coin to someone else")
// TODO: re-write this test to read from db
// 	for _, block := range bc.blocks {
// 		p := NewProofOfWork(block)
// 		assert.True(t, p.Validate(), "validation should be true")
// 	}
// }s

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
			got, err := NewBlockchain()
			assert.NoError(t, err, "no error")
			assert.IsType(t, &Blockchain{}, got, "if there's no error, this should be of type *Blockchain")
		})
	}
}
