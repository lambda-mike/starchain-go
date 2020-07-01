// proxy package provides BlockchainProxy struct which translates
// blockchain specific types to contracts and vice-versa.
package proxy

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/starchain/block"
	"github.com/starchain/blockchain"
	"github.com/starchain/contracts"
)

type BlockchainProxy struct {
	blockchain *blockchain.Blockchain
}

func New(blockchain *blockchain.Blockchain) BlockchainProxy {
	return BlockchainProxy{blockchain}
}

func (bp BlockchainProxy) RequestMessageOwnershipVerification(addr string) (string, error) {
	return bp.blockchain.RequestMessageOwnershipVerification(addr)
}

func (bp BlockchainProxy) GetBlockByHeight(h int) (contracts.Block, error) {
	var result contracts.Block
	block, err := bp.blockchain.GetBlockByHeight(h)
	if err == nil {
		result = MapBlockToContract(block)
	}
	return result, err
}

func (bp BlockchainProxy) GetBlockByHash(h string) (contracts.Block, error) {
	var (
		result contracts.Block
		hash   [sha256.Size]byte
	)
	buf, err := hex.DecodeString(h)
	if err != nil {
		return result, errors.New(fmt.Sprintf("Error occurred when decoding string: %s", err))
	}
	if len(buf) != sha256.Size {
		return result, errors.New(fmt.Sprintf("Hash must be exactly 32 char long, got: %v", len(buf)))
	}
	for i, _ := range buf {
		hash[i] = buf[i]
	}
	block, err := bp.blockchain.GetBlockByHash(hash)
	if err == nil {
		result = MapBlockToContract(block)
	}
	return result, err
}

func (bp BlockchainProxy) GetStarsByWalletAddress(addr string) []string {
	return bp.blockchain.GetStarsByWalletAddress(addr)
}

func (bp BlockchainProxy) SubmitStar(star contracts.StarData) (contracts.Block, error) {
	var req blockchain.StarRequest
	req.Addr = star.Address
	req.Msg = star.Message
	req.StarData = star.Data
	req.Sig = star.Signature
	b, err := bp.blockchain.SubmitStar(req)
	if err != nil {
		return contracts.Block{}, err
	}
	return MapBlockToContract(b), nil
}

func MapBlockToContract(block *block.Block) contracts.Block {
	var result contracts.Block
	result.Body = string(block.GetData())
	// TODO reuse utils.HashToStr fn
	result.Hash = fmt.Sprintf("%x", block.GetHash())
	result.Height = block.GetHeight()
	result.Owner = block.GetOwner()
	result.PreviousBlockHash = fmt.Sprintf("%x", block.GetPrevHash())
	result.Time = block.GetTimestamp()
	return result
}
