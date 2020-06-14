package block

import (
	"encoding/hex"
	"testing"
)

func TestNewBlock(t *testing.T) {
	t.Log("Block")
	{
		t.Log("TestNewBlock")
		{
			data := []byte("\"This is JSON string\"")
			t.Logf("\tGiven string data (%s)", data)
			{
				block := NewBlock(data)
				if block == nil {
					t.Errorf("\t\tShould return new Block, got: nil")
				}
				t.Log("\t\tShould return new Block")

				expected := make([]byte, hex.EncodedLen(len(data)))
				hex.Encode(expected, data)
				actual := block.data
				if len(actual) != len(expected) {
					t.Errorf("\t\tShould save data of proper length (%v), got: %v", len(expected), len(actual))
				}
				t.Logf("\t\tShould save data of proper length (%v)", len(expected))

				for i, b := range actual {
					if b != expected[i] {
						t.Errorf("\t\tShould contain the same byte at index: %v, got: %v ", i, b)
					}
				}
				t.Log("\t\tShould save data as hex")
			}
		}
	}
}
