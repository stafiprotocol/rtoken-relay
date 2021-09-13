package utils

import (
	"golang.org/x/crypto/blake2b"
)

func BlakeTwo256(dest []byte) [32]byte {
	return blake2b.Sum256(dest)
}
