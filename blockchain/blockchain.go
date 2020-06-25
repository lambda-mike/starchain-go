// This package provides a Blockchain structure.
// It depends on block package
package blockchain

// TODO add logging

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/starchain/block"
	"regexp"
	"strconv"
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

const FIVE_MIN int64 = 5 * 60

func (b BlockchainClock) GetTime() int64 {
	return time.Now().Unix()
}

var (
	EmptyAddrErr       = errors.New("Address is empty")
	EmptyMsgErr        = errors.New("Message is empty")
	EmptySigErr        = errors.New("Signature is empty")
	WrongTSErr         = errors.New("Message is not within allowed time range")
	MsgSigMistmatchErr = errors.New("Message does not match the signature")
)

// Factory function returning new Blockchain
func New(clock Clock) *Blockchain {
	var (
		blockchain Blockchain
	)
	owner := ""
	data := []byte("Genesis Gopher Block")
	blockchain.clock = clock
	blockchain.AddBlock(owner, data)
	return &blockchain
}

func (b *Blockchain) GetChainHeight() int {
	b.mutex.RLock()
	height := len(b.chain)
	b.mutex.RUnlock()
	return height
}

func (b *Blockchain) RequestMessageOwnershipVerification(addr string) (string, error) {
	if addr == "" {
		return "", EmptyAddrErr
	}
	b.mutex.RLock()
	ts := b.clock.GetTime()
	b.mutex.RUnlock()
	return fmt.Sprintf("%s:%d:starRegistry", addr, ts), nil
}

func (b *Blockchain) IsMessageOutdated(addr string, msg string) (bool, error) {
	regex := regexp.MustCompile(fmt.Sprintf("%s:(\\d{10,}):starRegistry", addr))
	if chunks := regex.FindStringSubmatch(msg); len(chunks) != 2 {
		return false, errors.New(fmt.Sprintf("Message %s is mlaformed", msg))
	} else if ts, err := strconv.ParseInt(chunks[1], 10, 64); err != nil {
		return false, errors.New(fmt.Sprintf("Chunk %v is not a number", chunks[1]))
	} else {
		b.mutex.RLock()
		now := b.clock.GetTime()
		b.mutex.RUnlock()
		duration := now - ts
		if duration < 0 {
			return true, WrongTSErr
		}
		return duration >= FIVE_MIN, nil
	}
	return true, errors.New("Undetermined result for timestamp checking")
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
	isOutdated, err := b.IsMessageOutdated(req.Addr, req.Msg)
	if err != nil {
		return nil, err
	}
	if isOutdated {
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

func (b *Blockchain) GetBlockByHash(hash [sha256.Size]byte) (*block.Block, error) {
	for _, b := range b.chain {
		if b.GetHash() == hash {
			return b, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Block %x not found", hash))
}

func (b *Blockchain) GetBlockByHeight(height int) (*block.Block, error) {
	if height < 0 || height >= len(b.chain) {
		return nil, errors.New(fmt.Sprintf("Invalid height: %v", height))
	}
	for _, block := range b.chain {
		if block.GetHeight() == height {
			return block, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Block at pos %v not found", height))
}

func (b *Blockchain) GetStarsByWalletAddress(addr string) []*block.Block {
	var blocks []*block.Block
	// Ommit genesis block - it has no owner
	for _, block := range b.chain[1:] {
		if block.GetOwner() == addr {
			blocks = append(blocks, block)
		}
	}
	return blocks
}

func (b *Blockchain) ValidateChain() bool {
	return false
}
