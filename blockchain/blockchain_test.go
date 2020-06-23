package blockchain

import (
	"fmt"
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
			// TODO prevHash?
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
