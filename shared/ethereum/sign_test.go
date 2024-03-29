package ethereum

import (
	"crypto/ecdsa"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestSignature(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("d50aafe92fe654904fdd6edfe714b2c3762c38a2157c377d39619a67112f4664")
	if err != nil {
		t.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatal("error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	t.Log("publicKeyBytes=", hexutil.Encode(publicKeyBytes))

	ethAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	t.Log("ethAddressByte", ethAddress.Bytes())
	t.Log("ethAddress", ethAddress.Hex())

	data, _ := hexutil.Decode("0x9c4189297ad2140c85861f64656d1d1318994599130d98b75ff094176d2ca31e")
	//data := []byte("hello")
	hash := crypto.Keccak256Hash(data)
	t.Log(hash.Hex())

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("signature", hexutil.Encode(signature))

	//sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log("sigPublicKey=", hexutil.Encode(sigPublicKey))
	//
	//matches := bytes.Equal(sigPublicKey, publicKeyBytes)
	//t.Log(matches) // true
	//
	//sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signature)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)
	//matches = bytes.Equal(sigPublicKeyBytes, publicKeyBytes)
	//t.Log(matches)
}
