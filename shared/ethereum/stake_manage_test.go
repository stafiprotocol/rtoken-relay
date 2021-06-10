package ethereum

import (
	"github.com/stafiprotocol/rtoken-relay/bindings/StakeManager"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

var (
	StakeManagerContract = common.HexToAddress("0x5e3Ef299fDDf15eAa0432E6e66473ace8c13D908")
	MainnetEndPoint      = "wss://mainnet.infura.io/ws/v3/86f8d5ba0d524274bce7780a83dbc0a4"
)

func TestSignerToValidatorId(t *testing.T) {
	client, err := NewSimpleClient(MainnetEndPoint)
	if err != nil {
		t.Fatal(err)
	}

	instance, err := StakeManager.NewStakeManager(StakeManagerContract, client)
	if err != nil {
		t.Fatal(err)
	}

	_ = instance

	signer := common.HexToAddress("0x49a7499e26f6311145934163d3aa286aea93fbe9")
	validatorId, err := instance.SignerToValidator(nil, signer)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(validatorId)
}
