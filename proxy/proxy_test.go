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

func TestGetBlockByHeight(t *testing.T) {
	t.Log("TestGetBlockByHeight")
	{
		bchain := blockchain.New(clock)
		proxy := New(bchain)
		h := 0
		t.Log("\tGiven a proper block height argument", h)
		{
			block, err := proxy.GetBlockByHeight(h)
			if err != nil {
				t.Fatal("\t\tShould return block without err, got err: ", err)
			}
			t.Log("\t\tShould return block without error")
			if block.Height != h {
				t.Fatal("\t\tShould return block with height", h, "got:", block.Height)
			}
			t.Log("\t\tShould return block with correct height")
			expectedBody := "Genesis Gopher Block"
			if block.Body != expectedBody {
				t.Fatal("\t\tShould return block with correct data, got:", block.Body)
			}
			t.Log("\t\tShould return block with correct data")
			if block.Hash != "8a9a61241b4825dfa8884c04678899974ddfde55532a2fbadc07fc78472c8731" {
				t.Fatal("\t\tShould return block with correct hash, got:", block.Hash)
			}
			t.Log("\t\tShould return block with correct hash")
			if block.Owner != "" {
				t.Fatal("\t\tShould return block with correct owner, got:", block.Owner)
			}
			t.Log("\t\tShould return block with correct owner")
			if block.PreviousBlockHash != "0000000000000000000000000000000000000000000000000000000000000000" {
				t.Fatal("\t\tShould return block with correct PreviousBlockHash, got:", block.PreviousBlockHash)
			}
			t.Log("\t\tShould return block with correct PreviousBlockHash")
		}
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
