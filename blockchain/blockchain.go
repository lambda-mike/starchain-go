// This package provides a Blockchain structure.
// It depends on block package
package blockchain

import (
	"github.com/starchain/block"
)

// Blockchain struct consists of the slice of Blocks (the chain).
// It allows basic operations such as adding block, checking height,
// checking blockchain integrity, fetching owner's blocks,
// getting block by id
type Blockchain struct {
	chain []block.Block
	// TODO mutex
}

// Factory function returning new Blockchain
func New() *Blockchain {
	return &Blockchain{}
}
