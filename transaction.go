package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"

	"github.com/labstack/gommon/log"
)

// Transaction is a combination of inputs and outputs
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

//TODO: Change reward dynamically
const subsidy int = 210000

//IsCoinbase returns whether or not this is the coinbase transaction
func (t *Transaction) IsCoinbase() bool {
	return len(t.Vin) == 1 && len(t.Vin[0].Txid) == 0 && t.Vin[0].Vout == -1
}

// SetID sets ID of a transaction
func (t *Transaction) SetID() error {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(t)
	if err != nil {
		return err
	}
	hash = sha256.Sum256(encoded.Bytes())
	t.ID = hash[:]
	return nil
}

// TXInput is the first part of a transaction
type TXInput struct {
	//Txid is the transaction ID
	Txid []byte
	// Vout stores the index of the TXOutput
	Vout int
	// Provides data to be use in the output's ScriptPubKey
	ScriptSig string
}

// CanUnlockOutputWith checks if the input can be "unlocked" with the given address TODO: Implement keys
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// TXOutput is the second part of a transaction
type TXOutput struct {
	// Value stores the "coin"
	Value int
	// ScriptPubKey stores addresses, currently
	ScriptPubKey string
}

// CanBeUnlockedWith checks if the input can be "unlocked" with the given address TODO: Implement keys
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

// NewCoinbaseTX is the initial block transaction, needing no earlier transactions
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s", to)
	}

	txin := TXInput{
		Txid:      []byte{},
		Vout:      -1,
		ScriptSig: data,
	}

	txout := TXOutput{
		Value:        subsidy,
		ScriptPubKey: to,
	}

	tx := Transaction{
		ID:   nil,
		Vin:  []TXInput{txin},
		Vout: []TXOutput{txout},
	}

	return &tx
}

// NewUTXOTransaction creates a transaction between addresses
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) (*Transaction, error) {
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs, err := bc.FindSpendableOutputs(from, amount)
	if err != nil {
		return nil, err
	}

	if acc < amount {
		return nil, fmt.Errorf("Not enough funds")
	}

	// List inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			return nil, err
		}

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from}) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	log.Info("NewUTXOTransaction: %s", tx)

	return &tx, tx.SetID()
}
