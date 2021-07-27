package ethereum

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stafiprotocol/rtoken-relay/bindings/Multisig"
)

func TestIsOwner(t *testing.T) {
	client, err := NewSimpleClient(goerliEndPoint)
	if err != nil {
		t.Fatal(err)
	}

	multiContract := common.HexToAddress("0xfc42de640aa9759d460e1a11416eca95d25c5908")
	multi, err := Multisig.NewMultisig(multiContract, client)
	if err != nil {
		t.Fatal(err)
	}

	//owner1 := common.HexToAddress("0x57530B16cD9a10E4997D816DF883CD6c1131CD22")
	//t.Log(owner1.Hex())
	//ownerFlag, err := multi.IsOwner(nil, owner1)
	//if err != nil {
	//	t.Fatal(err)
	//}

	//assert.True(t, ownerFlag)

	owners, err := multi.GetOwners(nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, owner := range owners {
		t.Log(owner)
	}
}

//func TestIsNotOwner(t *testing.T) {
//	client, err := NewSimpleClient(kovanEndPoint)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	multi, err := Multisig.NewMultisig(kovanMultisigContract, client)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	owner1 := common.HexToAddress("0x6238602a308762558DC140105a4cCb5C920EEed9")
//	ownerFlag, err := multi.IsOwner(nil, owner1)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	assert.False(t, ownerFlag)
//}
//
//func TestChainId(t *testing.T) {
//	client, err := NewSimpleClient(kovanEndPoint)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	multi, err := Multisig.NewMultisig(kovanMultisigContract, client)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	id, err := multi.GetChainId(nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//	assert.Equal(t, 0, id.Cmp(kovanMultisigContractChainId))
//}
//
//func TestMessageToSign(t *testing.T) {
//	client, err := NewSimpleClient(kovanEndPoint)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	multi, err := Multisig.NewMultisig(kovanMultisigContract, client)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	to := common.HexToAddress("0x94797870643B082f080f1ed7d52b61a58582f613")
//	value := big.NewInt(0)
//	cd, _ := hexutil.Decode("0xa9059cbb000000000000000000000000ad0bf51f7fc89e262edbbdf53c260088b024d8570000000000000000000000000000000000000000000000000000000000000001")
//
//	msg, err := multi.MessageToSign(nil, to, value, cd)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	t.Log(hexutil.Encode(msg[:])) //0x56c4da6876c6e053aadc33244e575b03ee499a9739166536a4bc755799561b75
//}
