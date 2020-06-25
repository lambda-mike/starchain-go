package blockchain

import (
	"crypto/sha256"
	"fmt"
	"github.com/starchain/block"
	"testing"
	"time"
)

type BlockchainClockMock struct{}

func (b BlockchainClockMock) GetTime() int64 {
	return time.Date(2020, time.June, 14, 17, 46, 32, 0, time.UTC).Unix()
}

func TestNew(t *testing.T) {
	t.Log("TestNew")
	{
		t.Log("\tWhen called")
		{
			clock := BlockchainClockMock{}
			blockchain := New(clock)
			if blockchain == nil {
				t.Fatalf("\t\tShould return new Blockchain, got:\nnil")
			}
			if len(blockchain.chain) <= 0 {
				t.Fatalf("\t\tShould contain genesis block, but is empty")
			}
			if string(blockchain.chain[0].GetData()) != "Genesis Gopher Block" {
				t.Fatalf("\t\tGenesis block should contain proper phrase")
			}
			if blockchain.chain[0].GetOwner() != "" {
				t.Fatalf("\t\tGenesis block should not have an owner")
			}
			t.Log("\t\tShould return new Blockchain with Genesis Block inside")
		}
	}
}

func TestGetChainHeight(t *testing.T) {
	t.Log("TestGetChainHeight")
	{
		t.Log("\tGiven fresh blockchain")
		{
			clock := BlockchainClockMock{}
			blockchain := New(clock)
			height := blockchain.GetChainHeight()
			if height != 1 {
				t.Fatalf("\t\tShould return 1, got: %v", height)
			}
			t.Log("\t\tShould return 1")
		}
		// TODO add block, test height
		// TODO add and get blocks in parallel
	}
}

func TestRequestMessageOwnershipVerification(t *testing.T) {
	t.Log("RequestMessageOwnershipVerification")
	{
		t.Log("\tGiven correct wallet address")
		{
			clock := BlockchainClockMock{}
			blockchain := New(clock)
			var addr = "1FzpnkhbAteDkU1wXDtd8kKizQhqWcsrWe"
			msg, err := blockchain.RequestMessageOwnershipVerification(addr)
			if err != nil {
				t.Fatal("\t\tShould return nil error, got: ", err)
			}
			expected := fmt.Sprintf("%s:%d:starRegistry", addr, 1592156792)
			if expected != msg {
				t.Fatal("\t\tShould return correct message, got: ", msg)
			}
			t.Log("\t\tShould return correct message")
		}
		t.Log("\tGiven empty wallet address")
		{
			clock := BlockchainClockMock{}
			blockchain := New(clock)
			var addr = ""
			msg, err := blockchain.RequestMessageOwnershipVerification(addr)
			if err != EmptyAddrErr {
				t.Fatal("\t\tShould return error, got: ", err)
			}
			if "" != msg {
				t.Fatal("\t\tShould return empty message, got: ", msg)
			}
			t.Log("\t\tShould return empty error")
		}
	}
}

func TestSubmitStar(t *testing.T) {
	t.Log("SubmitStar")
	{
		var (
			addr = "1FzpnkhbAteDkU1wXDtd8kKizQhqWcsrWe"
			msg  = fmt.Sprintf("%s:%d:starRegistry", addr, 1592156792-3*60)
			star = []byte("TODO star")
			sig  = "TODO sig"
			req  = StarRequest{addr, msg, star, sig}
		)
		t.Log("\tGiven correct params")
		{
			clock := BlockchainClockMock{}
			blockchain := New(clock)
			block, err := blockchain.SubmitStar(req)
			if err != nil {
				t.Fatal("\t\tShould return block without errors, got err: ", err)
			}
			if bData := block.GetData(); string(bData) != string(star) {
				t.Fatal("\t\tShould return block with original data, got: ", bData)
			}
			if bOwner := block.GetOwner(); bOwner != addr {
				t.Fatal("\t\tShould return block with correct owner, got: ", bOwner)
			}
			if bHeight := block.GetHeight(); bHeight != 1 {
				t.Fatal("\t\tShould return block with correct height, got: ", bHeight)
			}
			if bHeight := len(blockchain.chain); bHeight != 2 {
				t.Fatal("\t\tShould add block to the chain so it is 2 items long, got: ", bHeight)
			}
			prevH := blockchain.chain[1].GetPrevHash()
			genesisH := blockchain.chain[0].GetHash()
			if prevH != genesisH {
				t.Fatal("\t\tShould chain blocks together, got wrong prevHash: ",
					prevH,
					"\ninstead of: ",
					genesisH)
			}
			t.Log("\t\tShould return correct block and add it to the blockchain")
		}
	}
}

func TestIsMessageOutdated(t *testing.T) {
	t.Log("IsMessageOutdated")
	{
		var addr = "1FzpnkhbAteDkU1wXDtd8kKizQhqWcsrWe"
		t.Log("\tGiven correct message")
		{
			msg := fmt.Sprintf("%s:%d:starRegistry", addr, 1592156792-3*60)
			clock := BlockchainClockMock{}
			blockchain := New(clock)
			isOutdated, err := blockchain.IsMessageOutdated(addr, msg)
			if err != nil {
				t.Fatal("\t\tShould return false and nil err, got err: ", err)
			}
			if isOutdated {
				t.Fatal("\t\tShould return false, got: ", isOutdated)
			}
			t.Log("\t\tShould return false")
		}
		t.Log("\tGiven outdated message")
		{
			msg := fmt.Sprintf("%s:%d:starRegistry", addr, 1592156792-5*60)
			clock := BlockchainClockMock{}
			blockchain := New(clock)
			isOutdated, err := blockchain.IsMessageOutdated(addr, msg)
			if !isOutdated {
				t.Fatal("\t\tShould return true, got", isOutdated)
			}
			if err != nil {
				t.Fatal("\t\tShould return true and nil err, got err: ", err)
			}
			t.Log("\t\tShould return true")
		}
		t.Log("\tGiven malformed message")
		{
			// Message from the future
			msg := fmt.Sprintf("%s:%d:starRegistry", addr, 1592156792+5*60)
			clock := BlockchainClockMock{}
			blockchain := New(clock)
			_, err := blockchain.IsMessageOutdated(addr, msg)
			if err == nil {
				t.Fatal("\t\tShould return err, got nil")
			}
			t.Log("\t\tShould return err: ", err)
		}
	}
}

func TestGetBlockByHash(t *testing.T) {
	t.Log("GetBlockByHash")
	{
		var hash [sha256.Size]byte
		t.Log("\tGiven empty hash")
		{
			clock := BlockchainClockMock{}
			blockchain := New(clock)
			block, err := blockchain.GetBlockByHash(hash)

			if block != nil {
				t.Fatal("\t\tShould return nil, got: ")
			}
			t.Log("\t\tShould return nil")
			if err == nil {
				t.Fatal("\t\tShould return not nil err, got: ", err)
			}
			t.Log("\t\tShould return not nil err:", err)
		}
	}
	t.Log("\tGiven genesis block hash")
	{
		clock := BlockchainClockMock{}
		blockchain := New(clock)
		hash := blockchain.chain[0].GetHash()
		block, err := blockchain.GetBlockByHash(hash)

		if block == nil {
			t.Fatal("\t\tShould return genesis block, got: ", block)
		}
		t.Log("\t\tShould return genesis block")
		if err != nil {
			t.Fatal("\t\tShould not return err, got: ", err)
		}
		t.Log("\t\tShould return genesis block")
	}
	t.Log("\tGiven new block hash")
	{
		var (
			addr = "1FzpnkhbAteDkU1wXDtd8kKizQhqWcsrWe"
			msg  = fmt.Sprintf("%s:%d:starRegistry", addr, 1592156792-3*60)
			star = []byte("Brand new Star")
			sig  = "TODO sig"
			req  = StarRequest{addr, msg, star, sig}
		)
		clock := BlockchainClockMock{}
		blockchain := New(clock)
		blockchain.SubmitStar(req)
		hash := blockchain.chain[1].GetHash()
		block, err := blockchain.GetBlockByHash(hash)

		if block == nil || block != blockchain.chain[1] {
			t.Fatal("\t\tShould return new block, got: ", block)
		}
		t.Log("\t\tShould return new block")
		if err != nil {
			t.Fatal("\t\tShould not return err, got: ", err)
		}
		t.Log("\t\tShould return new block")
	}
}

func TestGetBlockByHeight(t *testing.T) {
	t.Log("GetBlockByHeight")
	{
		t.Log("\tGiven height -1")
		{
			clock := BlockchainClockMock{}
			blockchain := New(clock)
			block, err := blockchain.GetBlockByHeight(-1)

			if block != nil {
				t.Fatal("\t\tShould return nil, got: ")
			}
			t.Log("\t\tShould return nil")
			if err == nil {
				t.Fatal("\t\tShould return not nil err, got: ", err)
			}
			t.Log("\t\tShould return not nil err:", err)
		}
		t.Log("\tGiven height bigger than chain len")
		{
			clock := BlockchainClockMock{}
			blockchain := New(clock)
			block, err := blockchain.GetBlockByHeight(1)

			if block != nil {
				t.Fatal("\t\tShould return nil, got: ")
			}
			t.Log("\t\tShould return nil")
			if err == nil {
				t.Fatal("\t\tShould return not nil err, got: ", err)
			}
			t.Log("\t\tShould return not nil err:", err)
		}
	}
	t.Log("\tGiven height 0")
	{
		clock := BlockchainClockMock{}
		blockchain := New(clock)
		block, err := blockchain.GetBlockByHeight(0)

		if block == nil {
			t.Fatal("\t\tShould return genesis block, got: ", block)
		}
		t.Log("\t\tShould return genesis block")
		if err != nil {
			t.Fatal("\t\tShould not return err, got: ", err)
		}
		t.Log("\t\tShould return genesis block")
	}
	t.Log("\tGiven new block height")
	{
		var (
			addr = "1FzpnkhbAteDkU1wXDtd8kKizQhqWcsrWe"
			msg  = fmt.Sprintf("%s:%d:starRegistry", addr, 1592156792-3*60)
			star = []byte("Brand new Star")
			sig  = "TODO sig"
			req  = StarRequest{addr, msg, star, sig}
		)
		height := 1
		clock := BlockchainClockMock{}
		blockchain := New(clock)
		blockchain.SubmitStar(req)
		block, err := blockchain.GetBlockByHeight(height)

		if block == nil || block != blockchain.chain[height] {
			t.Fatal("\t\tShould return new block, got: ", block)
		}
		t.Log("\t\tShould return new block")
		if err != nil {
			t.Fatal("\t\tShould not return err, got: ", err)
		}
		t.Log("\t\tShould return new block")
	}
}

func TestGetStarsByWalletAddress(t *testing.T) {
	t.Log("GetStarsByWalletAddress")
	{
		t.Log("\tGiven empty address")
		{
			clock := BlockchainClockMock{}
			blockchain := New(clock)
			blocks := blockchain.GetStarsByWalletAddress("")

			if len(blocks) != 0 {
				t.Fatal("\t\tShould return empty array, got: ", blocks)
			}
			t.Log("\t\tShould return empty array")
		}
		t.Log("\tGiven blockchain with the new block")
		{
			t.Log("\tWhen proper address passed as param")
			{
				var (
					addr = "1FzpnkhbAteDkU1wXDtd8kKizQhqWcsrWe"
					msg  = fmt.Sprintf("%s:%d:starRegistry", addr, 1592156792-3*60)
					star = []byte("Brand new Star")
					sig  = "TODO sig"
					req  = StarRequest{addr, msg, star, sig}
				)
				clock := BlockchainClockMock{}
				blockchain := New(clock)
				blockchain.SubmitStar(req)
				blocks := blockchain.GetStarsByWalletAddress(addr)
				if len(blocks) != 1 {
					t.Fatal("\t\tShould return not empty array, got: ", blocks)
				}
				t.Log("\t\tShould return not empty array")
				if blocks[0].GetOwner() != addr || string(blocks[0].GetData()) != string(star) {
					t.Fatal("\t\tShould return proper block, got: ", blocks[0])
				}
				t.Log("\t\tShould return proper blocks")
			}
		}
		t.Log("\tGiven blockchain with new blocks belonging to different owners")
		{
			t.Log("\tWhen second owner's address is passed")
			{
				var (
					addr1 = "1FzpnkhbAteDkU1wXDtd8kKizQhqWcsrWe"
					addr2 = "1FzpnkhbAteDkU2wXDtd8kKizQhqWcsrWx"
					msg1  = fmt.Sprintf("%s:%d:starRegistry", addr1, 1592156792-3*60)
					msg2  = fmt.Sprintf("%s:%d:starRegistry", addr2, 1592156792-2*60)
					star1 = []byte("Brand new Star 1")
					star2 = []byte("Brand new Star 2")
					sig   = "TODO sig"
					req1  = StarRequest{addr1, msg1, star1, sig}
					req2  = StarRequest{addr2, msg2, star2, sig}
				)
				clock := BlockchainClockMock{}
				blockchain := New(clock)
				blockchain.SubmitStar(req1)
				blockchain.SubmitStar(req2)
				blocks := blockchain.GetStarsByWalletAddress(addr2)
				if len(blocks) != 1 {
					t.Fatal("\t\tShould return not empty array, got: ", blocks)
				}
				t.Log("\t\tShould return not empty array")
				if blocks[0].GetOwner() != addr2 || string(blocks[0].GetData()) != string(star2) {
					t.Fatal("\t\tShould return proper block, got: ", blocks[0])
				}
				t.Log("\t\tShould return proper blocks")
			}
		}
	}
}

func TestValidateChain(t *testing.T) {
	t.Log("ValidateChain")
	{
		t.Log("Given 1 block chain")
		{
			t.Log("\tWhen hash is valid")
			{
				clock := BlockchainClockMock{}
				blockchain := New(clock)
				isValid := blockchain.ValidateChain()
				if !isValid {
					t.Fatal("\t\tShould return true", isValid)
				}
				t.Log("\t\tShould return true")
			}
			t.Log("\tWhen hash is invalid")
			{
				var (
					data  []byte = []byte("\"This is JSON string\"")
					ts    int64  = time.Date(2020, time.June, 14, 17, 46, 32, 0, time.UTC).Unix()
					h     int    = 7
					owner string = "1FzpnkhbAteDkU1wXDtd8kKizQhqWcsrWe"
				)
				clock := BlockchainClockMock{}
				blockchain := New(clock)
				blockchain.chain[0] = block.New(ts, h, owner, nil, data)
				isValid := blockchain.ValidateChain()
				if isValid {
					t.Fatal("\t\tShould return false", isValid)
				}
				t.Log("\t\tShould return false")
			}
		}
		t.Log("Given 2 block chain")
		{
			t.Log("\tWhen second block prev hash is valid")
			{
				var (
					addr1 = "1FzpnkhbAteDkU1wXDtd8kKizQhqWcsrWe"
					msg1  = fmt.Sprintf("%s:%d:starRegistry", addr1, 1592156792-3*60)
					star1 = []byte("Brand new Star 1")
					sig   = "TODO sig"
					req1  = StarRequest{addr1, msg1, star1, sig}
				)
				clock := BlockchainClockMock{}
				blockchain := New(clock)
				blockchain.SubmitStar(req1)
				isValid := blockchain.ValidateChain()
				if !isValid {
					t.Fatal("\t\tShould return true, got: ", isValid)
				}
				t.Log("\t\tShould return true")
			}
			t.Log("\tWhen the first block is modified")
			{
				var (
					addr1 = "1FzpnkhbAteDkU1wXDtd8kKizQhqWcsrWe"
					msg1  = fmt.Sprintf("%s:%d:starRegistry", addr1, 1592156792-3*60)
					star1 = []byte("Brand new Star 1")
					sig   = "TODO sig"
					req1  = StarRequest{addr1, msg1, star1, sig}
				)
				var (
					data  []byte = []byte("\"This is JSON string\"")
					ts    int64  = time.Date(2020, time.June, 14, 17, 46, 32, 0, time.UTC).Unix()
					h     int    = 7
					owner string = "1FzpnkhbAteDkU1wXDtd8kKizQhqWcsrWe"
				)
				clock := BlockchainClockMock{}
				blockchain := New(clock)
				blockchain.SubmitStar(req1)
				blockchain.chain[0] = block.New(ts, h, owner, nil, data)
				isValid := blockchain.ValidateChain()
				if isValid {
					t.Fatal("\t\tShould return false, got: ", isValid)
				}
				t.Log("\t\tShould return false")
			}
		}
	}
}
