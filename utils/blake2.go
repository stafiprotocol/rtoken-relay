package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"golang.org/x/crypto/blake2b"
)

func Blake2Hash(num types.U32) (types.Hash, error) {
	defaultHash := types.NewHash([]byte(""))

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, num)
	if err != nil {
		return defaultHash, fmt.Errorf("binary.Write failed: %s", err)
	}

	h, err := blake2b.New256(buf.Bytes())
	if err != nil {

	}

	return types.NewHash(h.Sum([]byte(""))), nil
}
