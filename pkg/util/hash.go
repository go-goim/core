package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashString(str string) string {
	return hex.EncodeToString(Sha256([]byte(str)))
}

func Sha256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}
