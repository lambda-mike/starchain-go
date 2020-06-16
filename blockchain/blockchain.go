// This package provides a Blockchain structure.
// It depends on block package
package blockchain

import (
	"crypto/sha256"
	"github.com/starchain/block"
	"time"
)

// Blockchain struct consists of the slice of Blocks (the chain).
// It allows basic operations such as adding block, checking height,
// checking blockchain integrity, fetching owner's blocks,
// getting block by id
type Blockchain struct {
	chain []*block.Block
	// TODO mutex
}

// Factory function returning new Blockchain
func New() *Blockchain {
	var (
		blockchain                    = Blockchain{[]*block.Block{}}
		prevHash   *[sha256.Size]byte = nil
	)
	ts := time.Now().Unix()
	height := int64(len(blockchain.chain))
	owner := ""
	data := []byte("Genesis Gopher Block")
	genesisBlock := block.New(ts, height, owner, prevHash, data)
	blockchain.chain = append(blockchain.chain, genesisBlock)
	return &blockchain
}