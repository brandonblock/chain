package chain

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestBlockchain_AddBlock(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("send 1 coin to someone")
	bc.AddBlock("send 1 coin to someone else")

	assert.Equal(t, bc.blocks[0].Hash, bc.blocks[1].PrevBlockHash, "the previous hash should match the... previous hash")
}
