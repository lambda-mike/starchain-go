// This package provides a Blockchain structure.
// It depends on block package
package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/starchain/block"
	"sync"
	"time"
)

// Blockchain struct consists of the slice of Blocks (the chain).
// It allows basic operations such as adding block, checking height,
// checking blockchain integrity, fetching owner's blocks,
// getting block by id
type Blockchain struct {
	chain []*block.Block
	mutex sync.RWMutex
}

var (
	EmptyAddrErr = errors.New("Address is empty")
)

// Factory function returning new Blockchain
func New() *Blockchain {
	var (
		blockchain Blockchain
		prevHash   *[sha256.Size]byte = nil
	)
	ts := time.Now().Unix()
	height := len(blockchain.chain)
	owner := ""
	data := []byte("Genesis Gopher Block")
	genesisBlock := block.New(ts, height, owner, prevHash, data)
	blockchain.chain = append(blockchain.chain, genesisBlock)
	return &blockchain
}

func (b *Blockchain) GetChainHeight() int {
	b.mutex.RLock()
	height := len(b.chain)
	b.mutex.RUnlock()
	return height
}

// TODO ValidateChain

func RequestMessageOwnershipVerification(addr string) (string, error) {
	if addr == "" {
		return "", EmptyAddrErr
	}
	ts := time.Now().Unix()
	return fmt.Sprintf("%s:%d:starRegistry", addr, ts), nil
}

// TODO SubmitStar
// TODO GetBlockByHash
// TODO GetBlockByHeight
// TODO GetStarsByWalletAddress
