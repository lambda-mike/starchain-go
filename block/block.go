package block

import (
	"encoding/hex"
)

type Block struct {
	data []byte
}

// TODO add timestamp
// TODO add height
// TODO add hash
// TODO add prevHash
// TODO add owner
func NewBlock(data []byte) *Block {
	dataHex := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(dataHex, data)
	return &Block{dataHex}
}

func (b *Block) GetData() []byte {
	decoded := make([]byte, hex.DecodedLen(len(b.data)))
	hex.Decode(decoded, b.data)
	return decoded
}
