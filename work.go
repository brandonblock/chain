package chain

import "math/big"

// hashing difficulty
const targetBits = 24

// ProofOfWork store holds values used for demonstrating work
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork generates a new PoW record
func NewProofOfWork(b *Block) *ProofOfWork {

	// create target length (shifting the register left based on targetBits)
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	return &ProofOfWork{b, target}
}
