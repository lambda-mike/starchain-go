package blockchain

import (
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
