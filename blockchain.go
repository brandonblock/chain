package chain

// Blockchain holds a historical record of blocks
type Blockchain struct {
	blocks []*Block
}

// NewBlockchain returns a new Blockchain with initial genesis block
func NewBlockchain() *Blockchain {
	return &Blockchain{
		blocks: []*Block{NewGenesisBlock()},
	}
}

// AddBlock adds a block to the record
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}
