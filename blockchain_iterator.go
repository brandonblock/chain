package chain

import "github.com/boltdb/bolt"

// BlockchainIterator is a struct to easily handle the reading of blocks from the chain one-by-one
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}
