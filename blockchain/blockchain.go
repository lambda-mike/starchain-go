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
	EmptyAddrErr       = errors.New("Address is empty")
	EmptyMsgErr        = errors.New("Message is empty")
	EmptySigErr        = errors.New("Signature is empty")
	WrongTSErr         = errors.New("Message is outdated - wrong timestamp")
	MsgSigMistmatchErr = errors.New("Message does not match the signature")
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
	// TODO reues add method?
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

func (b *Blockchain) IsMessageOutdated(msg string) bool {
	// TODO
	return true
}

func (b *Blockchain) AddBlock(owner string, starData []byte) *block.Block {
	var prevHash [sha256.Size]byte
	b.mutex.Lock()
	ts := b.clock.GetTime()
	height := len(b.chain)
	if height > 0 {
		prevHash = b.chain[height-1].GetHash()
	}
	newBlock := block.New(ts, height, owner, &prevHash, starData)
	b.chain = append(b.chain, newBlock)
	b.mutex.Unlock()
	return newBlock
}

func (b *Blockchain) SubmitStar(req StarRequest) (*block.Block, error) {
	if req.Addr == "" {
		return nil, EmptyAddrErr
	}
	if req.Msg == "" {
		return nil, EmptyMsgErr
	}
	if req.Sig == "" {
		return nil, EmptySigErr
	}
	if b.IsMessageOutdated(req.Msg) {
		return nil, WrongTSErr
	}
	if !VerifyMessage(req) {
		return nil, MsgSigMistmatchErr
	}
	return b.AddBlock(req.Addr, req.StarData), nil
}

func VerifyMessage(req StarRequest) bool {
	// TODO verify msg based on the signature
	return true
}

// TODO GetBlockByHash
// TODO GetBlockByHeight
// TODO GetStarsByWalletAddress
