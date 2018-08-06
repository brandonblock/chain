package chain

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
			got, err := NewBlockchain()
			assert.NoError(t, err, "no error")
			assert.IsType(t, &Blockchain{}, got, "if there's no error, this should be of type *Blockchain")
		})
	}
	os.Remove(dbFile)

}

func TestBlockchain_AddBlock(t *testing.T) {
	bc, _ := NewBlockchain()

	err := bc.AddBlock("send 1 coin to someone")
	assert.NoError(t, err, "this should persist a transaction")
	err = bc.AddBlock("send 1 coin to someone else")
	assert.NoError(t, err, "this should persist another transaction")
	os.Remove(dbFile)

	// TODO: re-write this test to read from db
	// for _, block := range bc.blocks {
	// 	p := NewProofOfWork(block)
	// 	assert.True(t, p.Validate(), "validation should be true")
	// }
}

func BenchmarkBlockchain_AddBlock(b *testing.B) {
	bc, _ := NewBlockchain()
	for n := 0; n < b.N; n++ {
		bc.AddBlock("send 1 coin to someone")
	}
	os.Remove(dbFile)
}
