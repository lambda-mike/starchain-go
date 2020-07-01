// proxy package provides BlockchainProxy struct which translates
// blockchain specific types to contracts and vice-versa.
package proxy

import (
	"errors"
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
	// TODO
	var block contracts.Block
	switch h {
	default:
		return block, errors.New("TODO")
	}
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
