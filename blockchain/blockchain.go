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
	clock Clock
}

// StarRequest struct contains all data requiered to create a new star
// in the blockchain.
type StarRequest struct {
	Addr     string
	Msg      string
	StarData []byte
	Sig      string
}

type Clock interface {
	GetTime() int64
}

type BlockchainClock struct{}

func (b BlockchainClock) GetTime() int64 {
	return time.Now().Unix()
}

var (
	EmptyAddrErr = errors.New("Address is empty")
)

// Factory function returning new Blockchain
func New(clock Clock) *Blockchain {
	var (
		blockchain Blockchain
		prevHash   *[sha256.Size]byte = nil
	)
	ts := clock.GetTime()
	height := len(blockchain.chain)
	owner := ""
	data := []byte("Genesis Gopher Block")
	genesisBlock := block.New(ts, height, owner, prevHash, data)
	blockchain.chain = append(blockchain.chain, genesisBlock)
	blockchain.clock = clock
	return &blockchain
}

func (b *Blockchain) GetChainHeight() int {
	b.mutex.RLock()
	height := len(b.chain)
	b.mutex.RUnlock()
	return height
}

// TODO ValidateChain

func (b *Blockchain) RequestMessageOwnershipVerification(addr string) (string, error) {
	if addr == "" {
		return "", EmptyAddrErr
	}
	b.mutex.RLock()
	ts := b.clock.GetTime()
	b.mutex.RUnlock()
	return fmt.Sprintf("%s:%d:starRegistry", addr, ts), nil
}

func (b *Blockchain) SubmitStar(req StarRequest) (*block.Block, error) {
	return nil, nil
}

// TODO GetBlockByHash
// TODO GetBlockByHeight
// TODO GetStarsByWalletAddress
