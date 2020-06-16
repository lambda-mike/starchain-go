package blockchain

import (
	"fmt"
	"regexp"
	"testing"
)

func TestNew(t *testing.T) {
	t.Log("TestNew")
	{
		t.Log("\tWhen called")
		{
			blockchain := New()
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
			blockchain := New()
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
			var addr = "1FzpnkhbAteDkU1wXDtd8kKizQhqWcsrWe"
			msg, err := RequestMessageOwnershipVerification(addr)
			if err != nil {
				t.Fatal("\t\tShould return nil error, got: ", err)
			}
			regex := fmt.Sprintf("%s:\\d{10,}:starRegistry", addr)
			if matched, _ := regexp.MatchString(regex, msg); !matched {
				t.Fatal("\t\tShould return correct message, got: ", msg)
			}
			t.Log("\t\tShould return correct message")
		}
		t.Log("\tGiven empty wallet address")
		{
			var addr = ""
			msg, err := RequestMessageOwnershipVerification(addr)
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
			msg  = "TODO msg"
			star = []byte("TODO star")
			sig  = "TODO sig"
			req  = StarRequest{addr, msg, star, sig}
		)
		t.Log("\tGiven correct params")
		{
			blockchain := New()
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
			genesisH := blockchain.chain[0].CalculateHash()
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
