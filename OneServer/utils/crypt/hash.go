package crypt

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
