package main

import (
	"fmt"
	"os"

	"github.com/boltdb/bolt"
	"github.com/labstack/gommon/log"
)

const blocksBucket = "blocks"
const dbFile = "blockchain.db"
const genesisCoinbaseData = "You will never live if you are looking for the meaning of life."

// Blockchain holds a historical record of blocks
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// NewBlockchain returns a new Blockchain with initial genesis block
func NewBlockchain(address string) (*Blockchain, error) {

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	}); err != nil {
		log.Error(err)
		return nil, err
	}

	bc := Blockchain{tip, db}

	return &bc, nil
}

// CreateBlockchain creates a new blockchain DB
func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)

		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		serialized, err := genesis.Serialize()
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.Hash, serialized)
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

// AddBlock adds a block to the record
func (bc *Blockchain) AddBlock(transactions []*Transaction) (err error) {
	var lastHash []byte

	log.Infof("adding block data: %s", transactions)

	// read-only db transaction to get tip
	if err = bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	}); err != nil {
		return
	}

	// mine and persist new block
	newBlock := NewBlock(transactions, lastHash)
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

// Iterator returns a BlockchainIterator with the Blockchain values
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.db}
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
