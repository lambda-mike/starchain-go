package proxy

import (
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
	addr  string          = "1FzpnkhbAteDkU1wXDtd8kKizQhqWcsrWe"
	clock contracts.Clock = BlockchainClockMock{}
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
		h = 6
		t.Log("\tGiven a wrong block height argument", h)
		{
			_, err := proxy.GetBlockByHeight(h)
			if err == nil {
				t.Fatal("\t\tShould return err, got nil")
			}
			t.Log("\t\tShould return err:", err)
		}
	}
}

func TestGetBlockByHash(t *testing.T) {
	t.Log("TestGetBlockByHash")
	{
		bchain := blockchain.New(clock)
		proxy := New(bchain)
		hash := "8a9a61241b4825dfa8884c04678899974ddfde55532a2fbadc07fc78472c8731"
		t.Log("\tGiven a proper block hash argument", hash)
		{
			block, err := proxy.GetBlockByHash(hash)
			if err != nil {
				t.Fatal("\t\tShould return block without err, got err: ", err)
			}
			t.Log("\t\tShould return block without error")
			if block.Height != 0 {
				t.Fatal("\t\tShould return block with height", 0, "got:", block.Height)
			}
			t.Log("\t\tShould return block with correct height")
			expectedBody := "Genesis Gopher Block"
			if block.Body != expectedBody {
				t.Fatal("\t\tShould return block with correct data, got:", block.Body)
			}
			t.Log("\t\tShould return block with correct data")
			if block.Hash != hash {
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
		hash = "bada12bada12bada12bada12bada12bada12bada1bada1bada1bada122bada12"
		t.Log("\tGiven a wrong block hash argument", hash)
		{
			_, err := proxy.GetBlockByHash(hash)
			if err == nil {
				t.Fatal("\t\tShould return err, got nil")
			}
			t.Log("\t\tShould return err:", err)
		}
	}
}

func TestGetStarsByWalletAddress(t *testing.T) {
	t.Log("TestGetStarsByWalletAddress")
	{
		bchain := blockchain.New(clock)
		proxy := New(bchain)
		starsData := [][]byte{
			[]byte("Data 1"),
			[]byte("Data 2"),
			[]byte("Data 3"),
		}
		t.Log("\tGiven a proper block address argument and empty blockchain")
		{
			stars := proxy.GetStarsByWalletAddress(addr)
			if len(stars) != 0 {
				t.Fatal("\t\tShould return no stars, got:", len(stars))
			}
			t.Log("\t\tShould return no stars")
			bchain.AddBlock(addr, starsData[0])
			bchain.AddBlock("nope", starsData[1])
			bchain.AddBlock(addr, starsData[2])
			stars = proxy.GetStarsByWalletAddress(addr)
			if len(stars) != 2 {
				t.Fatal("\t\tShould return 2 stars, got:", len(stars))
			}
			if stars[0] != string(starsData[0]) {
				t.Fatalf("\t\tShould return correct star: %s, got: %s", starsData[0], stars[0])
			}
			if stars[1] != string(starsData[2]) {
				t.Fatalf("\t\tShould return correct star: %s, got: %s", starsData[2], stars[1])
			}
			t.Log("\t\tShould return correct 2 stars")
		}
	}
}

func TestSubmitStar(t *testing.T) {
	t.Log("TestSubmitStar")
	{
		bchain := blockchain.New(clock)
		proxy := New(bchain)
		t.Log("\tGiven proper star data")
		{
			var star contracts.StarData
			star.Address = addr
			star.Message = addr + ":1592156792:starRegistry"
			star.Data = []byte("New Star")
			star.Signature = "Sig"
			block, err := proxy.SubmitStar(star)
			if err != nil {
				t.Fatal("\t\tShould return block without err, got err: ", err)
			}
			t.Log("\t\tShould return block without error")
			if block.Height != 1 {
				t.Fatal("\t\tShould return block with height", 1, "got:", block.Height)
			}
			t.Log("\t\tShould return block with correct height")
			if block.Body != string(star.Data) {
				t.Fatal("\t\tShould return block with correct data, got:", block.Body)
			}
			t.Log("\t\tShould return block with correct data")
			if block.Hash != "c02b0315cdb0461a4bb4945ad132f28084ae8bb20f4038bc1ed3e181097aed2f" {
				t.Fatal("\t\tShould return block with correct hash, got:", block.Hash)
			}
			t.Log("\t\tShould return block with correct hash")
			if block.Owner != addr {
				t.Fatal("\t\tShould return block with correct owner, got:", block.Owner)
			}
			t.Log("\t\tShould return block with correct owner")
			if block.PreviousBlockHash != "8a9a61241b4825dfa8884c04678899974ddfde55532a2fbadc07fc78472c8731" {
				t.Fatal("\t\tShould return block with correct PreviousBlockHash, got:", block.PreviousBlockHash)
			}
		}
		t.Log("\tGiven wrong message")
		{
			var star contracts.StarData
			star.Address = addr
			star.Message = addr + ":1:starRegistry"
			star.Data = []byte("New Star")
			_, err := proxy.SubmitStar(star)
			if err == nil {
				t.Fatal("\t\tShould return err, got nil")
			}
			t.Log("\t\tShould return error:", err)
		}
	}
}
