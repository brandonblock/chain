package chain

import "github.com/boltdb/bolt"

// BlockchainIterator is a struct to easily handle the reading of blocks from the chain one-by-one
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Next returns the next block on the chain from i.currentHash
func (i *BlockchainIterator) Next() (*Block, error) {
	var block *Block

	if err := i.db.View(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block, err = DeserializeBlock(encodedBlock)
		return

	}); err != nil {
		return nil, err
	}

	i.currentHash = block.PrevBlockHash
	return block, nil
}
