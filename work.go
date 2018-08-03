package chain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/big"

	"github.com/labstack/gommon/log"
)

// hashing difficulty TODO: move to env or something for configuration
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

func (p *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	log.Infof("Mining block containing %v", p.block.Data)
	// don't overflow the counter
	for nonce < math.MaxInt64 {
		// prepare
		data := p.prepareData(nonce)
		// hash
		hash = sha256.Sum256(data)
		// convert has to big integer
		hashInt.SetBytes(hash[:])

		// compare hash with target
		if hashInt.Cmp(p.target) == -1 {
			return nonce, hash[:]
		}
		nonce++
	}
	return nonce, hash[:]
}

// Validate compares the current as to the current nonce
func (p *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := p.prepareData(p.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(p.target) == -1

	return isValid
}

// prepareData collates data for work verification
func (p *ProofOfWork) prepareData(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			p.block.PrevBlockHash,
			p.block.Data,
			intToHex(p.block.Timestamp),
			intToHex(int64(targetBits)),
			intToHex(int64(nonce)),
		},
		[]byte{},
	)
}

// intToHex converts an int64 to a byte array
func intToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
