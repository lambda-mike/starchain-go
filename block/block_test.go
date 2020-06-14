package block

import (
	"encoding/hex"
	"testing"
)

func TestNewBlock(t *testing.T) {
	t.Log("TestNewBlock")
	{
		data := []byte("\"This is JSON string\"")
		t.Logf("\tGiven string data (%s)", data)
		{
			block := NewBlock(data)
			if block == nil {
				t.Fatalf("\t\tShould return new Block, got: nil")
			}
			t.Log("\t\tShould return new Block")

			expected := make([]byte, hex.EncodedLen(len(data)))
			hex.Encode(expected, data)
			actual := block.data
			if len(actual) != len(expected) {
				t.Fatalf("\t\tShould save data of proper length (%v), got: %v", len(expected), len(actual))
			}
			t.Logf("\t\tShould save data of proper length (%v)", len(expected))

			for i, b := range actual {
				if b != expected[i] {
					t.Fatalf("\t\tShould contain the same byte at index: %v, got: %v ", i, b)
				}
			}
			t.Log("\t\tShould save data as hex")
		}
	}
}

func TestGetData(t *testing.T) {
	t.Log("GetData")
	{
		data := []byte("\"This is random JSON string\"")
		t.Logf("\tGiven a block with data (%s)", data)
		{
			block := NewBlock(data)
			actual := block.GetData()

			if string(actual) != string(data) {
				t.Fatalf("\t\tShould return the same data, got: (%s)", actual)
			}
			t.Logf("\t\tShould return the same data.")
		}
	}
}

func TestValidate(t *testing.T) {
	t.Log("Validate")
	{
		t.Logf("\tGiven a block without data (nil)")
		{
			block := NewBlock(nil)
			isValid := block.Validate()
			if !isValid {
				t.Fatalf("\t\tShould return true, got: %v", isValid)
			}
			t.Logf("\t\tShould return true.")
		}
		data := []byte("\"This is original data\"")
		t.Logf("\tGiven a block with data (%s)", data)
		{
			block := NewBlock(data)
			isValid := block.Validate()
			if !isValid {
				t.Fatalf("\t\tShould return true, got: %v", isValid)
			}
			t.Logf("\t\tShould return true.")
		}
		t.Logf("\tGiven a block with data (%s)", data)
		{
			t.Log("\t\tWhen data was changed")
			{
				block := NewBlock(data)
				block.data = []byte("Not this time!")
				isValid := block.Validate()
				if isValid {
					t.Fatal("\t\t\tShould return false")
				}
				t.Log("\t\t\tShould return false")
			}
		}
	}
}
