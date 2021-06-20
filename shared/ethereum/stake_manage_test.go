package ethereum

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/rtoken-relay/bindings/StakeManager"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

var (
	mainnetStakeManagerContract = common.HexToAddress("0x5e3Ef299fDDf15eAa0432E6e66473ace8c13D908")
	MainnetEndPoint             = "wss://mainnet.infura.io/ws/v3/86f8d5ba0d524274bce7780a83dbc0a4"

	goerliStakeManagerContract = common.HexToAddress("0x4864d89DCE4e24b2eDF64735E014a7E4154bfA7A")
)

func TestSignerToValidatorId(t *testing.T) {
	client, err := NewSimpleClient(MainnetEndPoint)
	if err != nil {
		t.Fatal(err)
	}

	instance, err := StakeManager.NewStakeManager(mainnetStakeManagerContract, client)
	if err != nil {
		t.Fatal(err)
	}

	signer := common.HexToAddress("0x49a7499e26f6311145934163d3aa286aea93fbe9")
	validatorId, err := instance.SignerToValidator(nil, signer)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(validatorId) // 115
}

func TestSigners(t *testing.T) {
	client, err := NewSimpleClient(MainnetEndPoint)
	if err != nil {
		t.Fatal(err)
	}

	instance, err := StakeManager.NewStakeManager(mainnetStakeManagerContract, client)
	if err != nil {
		t.Fatal(err)
	}

	addr, err := instance.Signers(nil, big.NewInt(1))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(addr) // 0x0306B7d3095Ab008927166CD648a8ca7dBe53F05

	id, err := instance.SignerToValidator(nil, addr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id) // 111
	t.Log(hexutil.Encode(id.Bytes()))

	share, err := instance.Validators(nil, id)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(share.ContractAddress) // 0x59a82694a675377E010F18F598f5B9dBB83eD968
	t.Log(share.Status) // 1

	valFlag, err := instance.IsValidator(nil, id)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(valFlag)
}
