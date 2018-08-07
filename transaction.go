package main

import "fmt"

// Transaction is a combination of inputs and outputs
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

//TODO: Change reward dynamically
const subsidy int = 210000

// TXInput is the first part of a transaction
type TXInput struct {
	//Txid is the transaction ID
	Txid []byte
	// Vout stores the index of the TXOutput
	Vout int
	// Provides data to be use in the output's ScriptPubKey
	ScriptSig string
}

// TXOutput is the second part of a transaction
type TXOutput struct {
	// Value stores the "coin"
	Value int
	// ScriptPubKey stores addresses, currently
	ScriptPubKey string
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
