package chain

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestBlockchain_AddBlock(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("send 1 coin to someone")
	bc.AddBlock("send 1 coin to someone else")

	for _, block := range bc.blocks {
		p := NewProofOfWork(block)
		assert.True(t, p.Validate(), "validation should be true")
	}
}
