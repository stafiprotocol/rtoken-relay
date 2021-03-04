package utils

import (
	"encoding/json"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"golang.org/x/crypto/blake2b"
)

func Blake2Hash(dest interface{}) (types.Hash, error) {
	bz, err := json.Marshal(dest)
	if err != nil {
		return types.NewHash([]byte{}), err
	}

	h, err := blake2b.New256(bz)
	if err != nil {
		return types.NewHash([]byte{}), err
	}

	return types.NewHash(h.Sum([]byte(""))), nil
}
