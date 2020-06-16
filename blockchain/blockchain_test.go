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
			t.Log("\t\tShould return new Blockchain")
		}
	}
}
