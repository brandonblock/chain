package main

import (
	"bytes"
	"encoding/gob"
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
	nonce, hash := p.Prove()
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

// DeserializeBlock returns a block object from a serialized []byte
func DeserializeBlock(d []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)

	return &block, err
}

// Serialize the block into an encoding/gob type to prepare for persistence
func (b *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(b); err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}
