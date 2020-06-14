package block

import (
	"testing"
)

func TestNewBlock(t *testing.T) {
	t.Log("Block")
	{
		t.Log("TestNewBlock")
		{
			const data = "\"This is JSON string\""
			t.Logf("\tGiven string data (%v)", data)
			{
				block := NewBlock(data)
				if block == nil {
					t.Errorf("\t\tShould return new Block, got: nil")
				}
				if block.data != data {
					t.Errorf("\t\tNew block should contain data, got: %v", block.data)
				}
			}
		}
	}
}
