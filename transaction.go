package main

// Transaction is a combination of inputs and outputs
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
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

// TXOutput is the second part of a transaction
type TXOutput struct {
	// Value stores the "coin"
	Value int
	// ScriptPubKey stores addresses, currently
	ScriptPubKey string
}
