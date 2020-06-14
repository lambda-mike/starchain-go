package block

import (
	"encoding/hex"
)

type Block struct {
	data []byte
}

func NewBlock(data []byte) *Block {
	dataHex := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(dataHex, data)
	return &Block{dataHex}
}
