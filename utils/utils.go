package utils

import (
	"crypto/sha256"
	"fmt"
)

func HashToStr(hash [sha256.Size]byte) string {
	return fmt.Sprintf("%x", hash)
}
