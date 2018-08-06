package chain

import (
	"os"

	"github.com/boltdb/bolt"
)

const blocksBucket = "blocks"
const dbFile = "blockchain.db"

// Blockchain holds a historical record of blocks
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// NewBlockchain returns a new Blockchain with initial genesis block
func NewBlockchain() (*Blockchain, error) {
	dbFile := os.Getenv(dbFile)
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}
	var tip []byte

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				return err
			}
			serialized, err := genesis.Serialize()
			if err != nil {
				return err
			}
			if err = b.Put(genesis.Hash, serialized); err != nil {
				return err
			}
			if err = b.Put([]byte("l"), genesis.Hash); err != nil {
				return err
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	return &Blockchain{tip, db}, nil
}

// Update adds a transaction to the blockchain
func update(tx *bolt.Tx) (tip []byte, err error) {
	b := tx.Bucket([]byte(blocksBucket))

	if b == nil {
		genesis := NewGenesisBlock()
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			return nil, err
		}
		serialized, err := genesis.Serialize()
		if err != nil {
			return nil, err
		}
		if err = b.Put(genesis.Hash, serialized); err != nil {
			return nil, err
		}
		if err = b.Put([]byte("l"), genesis.Hash); err != nil {
			return nil, err
		}
		tip = genesis.Hash
	} else {
		tip = b.Get([]byte("l"))
	}

	return tip, nil
}

// AddBlock adds a block to the record
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}
