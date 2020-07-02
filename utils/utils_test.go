package utils

import (
	"crypto/sha256"
	"testing"
)

func TestHashToStr(t *testing.T) {
	t.Log("HashToStr")
	{
		t.Log("\tGiven an empty hash")
		{
			var hash [sha256.Size]byte
			str := HashToStr(hash)
			expected := "0000000000000000000000000000000000000000000000000000000000000000"
			if str != expected {
				t.Fatal("\t\tShould return correct string, got:", str)
			}
			t.Log("\t\tShould return correct string")
		}
		t.Log("\tGiven a non-empty hash")
		{
			var hash [sha256.Size]byte
			for i := 0; i < sha256.Size; i++ {
				hash[i] = byte(i + 64)
			}
			str := HashToStr(hash)
			expected := "404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"
			if str != expected {
				t.Fatal("\t\tShould return correct string, got:", str)
			}
			t.Log("\t\tShould return correct string")
		}
	}
}
