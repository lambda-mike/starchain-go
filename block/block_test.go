package block

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"
)

var (
	data  []byte            = []byte("\"This is JSON string\"")
	ts    int64             = time.Date(2020, time.June, 14, 17, 46, 32, 0, time.UTC).Unix()
	h     int64             = 0
	prevH [sha256.Size]byte = sha256.Sum256([]byte("Here goes proper hash of the block"))
	owner string            = "1FzpnkhbAteDkU1wXDtd8kKizQhqWcsrWe"
)

func TestNewBlock(t *testing.T) {
	t.Log("TestNewBlock")
	{
		t.Log("\tGiven correct paramerters: ", ts, h, owner, data)
		{
			block := NewBlock(ts, h, owner, &prevH, data)
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

			for i, b := range block.prevHash {
				if prevH[i] != b {
					t.Fatalf("\t\tShould contain the same byte for prevHash at index: %v, got: %v ", i, b)
				}
			}
			t.Log("\t\tShould have correct prevHash bytes")
		}
	}
}

func TestNewBlockNilPrevH(t *testing.T) {
	t.Log("TestNewBlock")
	{
		t.Log("\tGiven nil hash of prev block")
		{
			block := NewBlock(ts, h, owner, nil, data)
			if block.prevHash != nil {
				t.Fatal("\t\tShould construct valid Block, prevHash should be nil")
			}
			t.Log("\t\tShould construct valid Block")
		}
	}
}

func TestNewBlockBadTS(t *testing.T) {
	t.Log("TestNewBlock")
	{
		var badTS int64 = 0
		t.Log("\tGiven incorrect timestamp", badTS)
		{
			defer func() {
				err := recover()
				if err != nil {
					t.Log("\t\tShould panic", err)
					return
				}
				t.Fatal("\t\tShould panic but got nil instead")
			}()
			_ = NewBlock(badTS, h, owner, &prevH, data)
		}
	}
}

func TestNewBlockBadHeight(t *testing.T) {
	t.Log("TestNewBlock")
	{
		var badHeight int64 = -1
		t.Log("\tGiven incorrect height", badHeight)
		{
			defer func() {
				err := recover()
				if err != nil {
					t.Log("\t\tShould panic", err)
					return
				}
				t.Fatal("\t\tShould panic but got nil instead")
			}()
			_ = NewBlock(ts, badHeight, owner, &prevH, data)
		}
	}
}

func TestGetData(t *testing.T) {
	t.Log("GetData")
	{
		data := []byte("\"This is random JSON string\"")
		t.Logf("\tGiven a block with data (%s)", data)
		{
			block := NewBlock(ts, h, owner, &prevH, data)
			actual := block.GetData()

			if string(actual) != string(data) {
				t.Fatalf("\t\tShould return the same data, got: (%s)", actual)
			}
			t.Logf("\t\tShould return the same data")
		}
	}
}

func TestValidate(t *testing.T) {
	t.Log("Validate")
	{
		t.Logf("\tGiven a block without data (nil)")
		{
			block := NewBlock(ts, h, owner, &prevH, nil)
			isValid := block.Validate()
			if !isValid {
				t.Fatalf("\t\tShould return true, got: %v", isValid)
			}
			t.Logf("\t\tShould return true")
		}
		data := []byte("\"This is original data\"")
		t.Logf("\tGiven a block with data (%s)", data)
		{
			block := NewBlock(ts, h, owner, &prevH, data)
			isValid := block.Validate()
			if !isValid {
				t.Fatalf("\t\tShould return true, got: %v", isValid)
			}
			t.Logf("\t\tShould return true")
		}
		t.Logf("\tGiven a block with data (%s)", data)
		{
			t.Log("\t\tWhen data was changed")
			{
				block := NewBlock(ts, h, owner, &prevH, data)
				block.data = []byte(string(data) + " Not this time!")
				isValid := block.Validate()
				if isValid {
					t.Fatal("\t\t\tShould return false, but got true")
				}
				t.Log("\t\t\tShould return false")
			}
		}
		t.Logf("\tGiven a block with ts (%v)", ts)
		{
			t.Log("\t\tWhen ts was changed")
			{
				block := NewBlock(ts, h, owner, &prevH, data)
				block.ts = 1
				isValid := block.Validate()
				if isValid {
					t.Fatal("\t\t\tShould return false, but got true")
				}
				t.Log("\t\t\tShould return false")
			}
		}
		t.Logf("\tGiven a block with height (%v)", h)
		{
			t.Log("\t\tWhen height was changed")
			{
				block := NewBlock(ts, h, owner, &prevH, data)
				block.height = h + 11
				isValid := block.Validate()
				if isValid {
					t.Fatal("\t\t\tShould return false, but got true")
				}
				t.Log("\t\t\tShould return false")
			}
		}
		t.Logf("\tGiven a block with owner (%v)", owner)
		{
			t.Log("\t\tWhen owner was changed")
			{
				block := NewBlock(ts, h, owner, &prevH, data)
				block.owner += "boom!"
				isValid := block.Validate()
				if isValid {
					t.Fatal("\t\t\tShould return false, but got true")
				}
				t.Log("\t\t\tShould return false")
			}
		}
		t.Logf("\tGiven a block with prevHash (%s)", prevH)
		{
			t.Log("\t\tWhen prevHash was changed")
			{
				block := NewBlock(ts, h, owner, &prevH, data)
				block.prevHash[2] = 22
				isValid := block.Validate()
				if isValid {
					t.Fatal("\t\t\tShould return false, but got true")
				}
				t.Log("\t\t\tShould return false")
			}
		}
	}
}
