package proxy

import (
	"errors"
	"github.com/starchain/blockchain"
	"github.com/starchain/contracts"
	"testing"
	"time"
)

type BlockchainClockMock struct{}

func (b BlockchainClockMock) GetTime() int64 {
	return time.Date(2020, time.June, 14, 17, 46, 32, 0, time.UTC).Unix()
}

var (
	addr string = "1FzpnkhbAteDkU1wXDtd8kKizQhqWcsrWe"
	// TODO move Clock to contracts
	clock blockchain.Clock = BlockchainClockMock{}
)

func TestRequestMessageOwnershipVerification(t *testing.T) {
	t.Log("TestRequestMessageOwnershipVerification")
	{
		bchain := blockchain.New(clock)
		proxy := New(bchain)
		t.Log("\tGiven an address: ", addr)
		{
			msg, err := proxy.RequestMessageOwnershipVerification(addr)
			if err != nil {
				t.Fatal("\t\tShould return message without err, got err: ", err)
			}
			t.Log("\t\tShould return message without error")
			expectedMsg := "1FzpnkhbAteDkU1wXDtd8kKizQhqWcsrWe:1592156792:starRegistry"
			if msg != expectedMsg {
				t.Fatal("\t\tShould return correct message, got: ", msg)
			}
			t.Log("\t\tShould return correct message")
		}
	}
}

func (b BlockchainProxy) TestGetBlockByHeight(h int) (contracts.Block, error) {
	// TODO
	var block contracts.Block
	switch h {
	default:
		return block, errors.New("TODO")
	}
}

func (b BlockchainProxy) TestGetBlockByHash(h string) (contracts.Block, error) {
	// TODO
	var block contracts.Block
	switch h {
	default:
		return block, errors.New("TODO")
	}
}

func (b BlockchainProxy) TestGetStarsByWalletAddress(addr string) []string {
	// TODO
	var stars []string = make([]string, 0)
	return stars
}

func (b BlockchainProxy) TestSubmitStar(star contracts.StarData) (contracts.Block, error) {
	// TODO
	var block contracts.Block
	return block, nil
}
