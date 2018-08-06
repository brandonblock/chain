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
	// open db file
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}
	var tip []byte

	// create read-write boltdb transaction
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
func (bc *Blockchain) AddBlock(data string) (err error) {
	var lastHash []byte

	if err = bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	}); err != nil {
		return
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(blocksBucket))
		serialized, err := newBlock.Serialize()
		if err != nil {
			return
		}
		if err = b.Put(newBlock.Hash, serialized); err != nil {
			return
		}
		if err = b.Put([]byte("l"), newBlock.Hash); err != nil {
			return
		}
		bc.tip = newBlock.Hash
		return
	})

	return
}
