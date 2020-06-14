package block

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Block struct {
	data []byte
	hash [sha256.Size]byte
}

// TODO add timestamp
// TODO add height
// TODO add prevHash
// TODO add owner
func NewBlock(data []byte) *Block {
	var block Block
	if data != nil {
		dataHex := make([]byte, hex.EncodedLen(len(data)))
		hex.Encode(dataHex, data)
		block.data = dataHex
	}
	hash := block.CalculateHash()
	block.hash = hash
	return &block
}

// CalculateHash method calculates the sha256 hash of the block properties
// except the hash field and returns that value.
func (b *Block) CalculateHash() [sha256.Size]byte {
	// TODO other props!!
	data := ""
	if b.data != nil {
		data = fmt.Sprintf("%s", b.data)
	}
	blockFields := fmt.Sprintf("%s", data)
	return sha256.Sum256([]byte(blockFields))
}

// GetData method returns data stored inside block decoded from hex
func (b *Block) GetData() []byte {
	decoded := make([]byte, hex.DecodedLen(len(b.data)))
	hex.Decode(decoded, b.data)
	return decoded
}

// Validate method checks whether block was tampered with.
// It does so by calculating the hash of the block without hash field and
// comparing the result with the hash stored in that block.
func (b *Block) Validate() bool {
	return b.hash == b.CalculateHash()
}
