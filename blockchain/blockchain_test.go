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
