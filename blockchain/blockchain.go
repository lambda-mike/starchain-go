// This package provides a Blockchain structure.
// It depends on block package
package blockchain

// Blockchain struct consists of the slice of Blocks (the chain).
// It allows basic operations such as adding block, checking height,
// checking blockchain integrity, fetching owner's blocks,
// getting block by id
type Blockchain struct{}

// Factory function returning new Blockchain
func New() *Blockchain {
	return &Blockchain{}
}
