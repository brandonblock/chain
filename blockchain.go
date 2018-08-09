package main

import (
	"encoding/hex"
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
func CreateBlockchain(address string) (*Blockchain, error) {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)

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

		return nil
	})

	if err != nil {
		return nil, err
	}

	bc := Blockchain{tip, db}

	return &bc, nil
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

// FindUnspentTransactions finds outputs that have no input references
func (bc *Blockchain) FindUnspentTransactions(address string) ([]Transaction, error) {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block, err := bci.Next()
		if err != nil {
			return nil, err
		}

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			// Check every output on the chain
		Outputs:
			for outIdx, out := range tx.Vout {

				// if the output is already "claimed" by an input, we ignore
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				// We want outputs that can unlock with our given address
				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			// coinbase transactions don't unlock inputs so we can ignore
			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			return unspentTXs, nil
		}
	}
}

// FindUTXO returns unspent transactions for a given address
func (bc *Blockchain) FindUTXO(address string) ([]TXOutput, error) {
	var UTXOs []TXOutput
	unspentTransactions, err := bc.FindUnspentTransactions(address)
	if err != nil {
		return nil, err
	}

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs, nil
}

// FindSpendableOutputs iterates over unspent transactions and grabs their values to spend
func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int, error) {
	unspentOutputs := make(map[string][]int)
	unspentTXs, err := bc.FindUnspentTransactions(address)
	if err != nil {
		return 0, nil, err
	}
	accumulated := 0

	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		// Accumulate unspent outputs until they equal the spend amount, then return the spend amount and the remainder
		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					return accumulated, unspentOutputs, nil
				}
			}
		}
	}
	return accumulated, unspentOutputs, nil
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
