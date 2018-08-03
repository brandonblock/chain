package chain

import (
	"time"
)

// Block stores on-chain information
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// NewGenesisBlock creates this first block on a chain
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// NewBlock creates a new Block and sets its hash
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
	}

	//prove our work
	p := NewProofOfWork(&block)
	nonce, hash := p.Run()
	block.Hash = hash
	block.Nonce = nonce
	return &block
}
