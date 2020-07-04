package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/starchain/utils"
	"io"
	"log"
)

// Block struct represents single block in the blockchain.
// It consists of timestamp (ts), height, address of the owner's wallet,
// previous block hash, data encoded as []byte of hex values,
// and SHA256 hash of the block.
type Block struct {
	ts       int64
	height   int
	owner    string
	prevHash *[sha256.Size]byte
	data     []byte
	hash     [sha256.Size]byte
}

var (
	WrongTimeStampErr error = errors.New("Timestamp must be bigger than 0")
	NegativeHeightErr error = errors.New("Height must be greater than or equal 0")
)

// New fn creates a brand new Block.
// It panics when timestamp is less than or equal 0.
// It panics when height is negative.
func New(ts int64, height int, owner string, prevHash *[sha256.Size]byte, data []byte) *Block {
	if ts <= 0 {
		log.Panic(WrongTimeStampErr, ts)
	}
	if height < 0 {
		log.Panic(NegativeHeightErr, height)
	}
	var block Block
	block.ts = ts
	block.height = height
	block.owner = owner
	if prevHash == nil {
		block.prevHash = nil
	} else {
		block.prevHash = new([sha256.Size]byte)
		copy(block.prevHash[:], prevHash[:])
	}
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
	data := ""
	if b.data != nil {
		data = fmt.Sprintf("%s", b.data)
	}
	prevH := ""
	if b.prevHash != nil {
		prevH = utils.HashToStr(*b.prevHash)
	}
	blockFields := fmt.Sprintf("|%d|%d|%s|%s|%s|", b.ts, b.height, b.owner, prevH, data)
	return sha256.Sum256([]byte(blockFields))
}

// DecodeData method returns data stored inside a block decoded from hex
func (b *Block) DecodeData() []byte {
	decoded := make([]byte, hex.DecodedLen(len(b.data)))
	hex.Decode(decoded, b.data)
	return decoded
}

// GetData method returns data stored inside a block decoded in hex format
func (b *Block) GetData() []byte {
	buffer := bytes.NewBufferString("")
	n, err := io.Copy(buffer, bytes.NewBuffer(b.data))
	if err != nil || n != int64(len(b.data)) {
		log.Panic("ERR: GetData failed to return a copy of Block.data", err)
	}
	return buffer.Bytes()
}

// GetOwner method returns owner stored inside a block
func (b *Block) GetOwner() string {
	return b.owner
}

// GetHeight method returns height of the block
func (b *Block) GetHeight() int {
	return b.height
}

// GetPrevHash method returns prevHash field value
func (b *Block) GetPrevHash() [sha256.Size]byte {
	var hash [sha256.Size]byte
	copy(hash[:], b.prevHash[:])
	return hash
}

// GetHash method returns hash field value
func (b *Block) GetHash() [sha256.Size]byte {
	var h [sha256.Size]byte
	copy(h[:], b.hash[:])
	return h
}

// GetTimestamp
func (b *Block) GetTimestamp() int64 {
	return b.ts
}

// Validate method checks whether block was tampered with.
// It does so by calculating the hash of the block without hash field and
// comparing the result with the hash stored in that block.
func (b *Block) Validate() bool {
	return b.hash == b.CalculateHash()
}
