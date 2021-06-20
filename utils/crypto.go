package utils

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
)

func DecomposeSignature(sigs [][]byte) ([]uint8, [][32]byte, [][32]byte) {
	rs := make([][32]byte, 0)
	ss := make([][32]byte, 0)
	vs := make([]uint8, 0)

	for _, sig := range sigs {
		var r [32]byte
		var s [32]byte
		copy(r[:], sig[:32])
		copy(s[:], sig[32:64])
		rs = append(rs, r)
		ss = append(ss, s)
		vs = append(vs, sig[64:][0])
	}

	return vs, rs, ss
}

func PublicKeyFromKeypair(pair *secp256k1.Keypair) []byte {
	publicKey := pair.PrivateKey().Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	return crypto.FromECDSAPub(publicKeyECDSA)
}
