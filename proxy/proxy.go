// proxy package provides BlockchainProxy struct which translates
// blockchain specific types to contracts and vice-versa.
package proxy

import (
	"errors"
	"fmt"
	"github.com/starchain/blockchain"
	"github.com/starchain/contracts"
)

type BlockchainProxy struct {
	blockchain *blockchain.Blockchain
}

func New(blockchain *blockchain.Blockchain) BlockchainProxy {
	return BlockchainProxy{blockchain}
}

func (b BlockchainProxy) RequestMessageOwnershipVerification(addr string) (string, error) {
	return b.blockchain.RequestMessageOwnershipVerification(addr)
}

func (b BlockchainProxy) GetBlockByHeight(h int) (contracts.Block, error) {
	var result contracts.Block
	block, err := b.blockchain.GetBlockByHeight(h)
	if err == nil {
		result.Body = string(block.GetData())
		// TODO reuse utils.HashToStr fn
		result.Hash = fmt.Sprintf("%x", block.GetHash())
		result.Height = block.GetHeight()
		result.Owner = block.GetOwner()
		result.PreviousBlockHash = fmt.Sprintf("%x", block.GetPrevHash())
		result.Time = block.GetTimestamp()
	}
	return result, err
}

func (b BlockchainProxy) GetBlockByHash(h string) (contracts.Block, error) {
	// TODO
	var block contracts.Block
	switch h {
	default:
		return block, errors.New("TODO")
	}
}

func (b BlockchainProxy) GetStarsByWalletAddress(addr string) []string {
	// TODO
	var stars []string = make([]string, 0)
	return stars
}

func (b BlockchainProxy) SubmitStar(star contracts.StarData) (contracts.Block, error) {
	// TODO
	var block contracts.Block
	return block, nil
}
