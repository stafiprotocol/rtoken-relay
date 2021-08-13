package ethereum

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/rtoken-relay/bindings/MaticToken"
	"github.com/stafiprotocol/rtoken-relay/bindings/StakeManager"
	"github.com/stafiprotocol/rtoken-relay/bindings/ValidatorShare"
)

var (
	mainnetStakeManagerContract = common.HexToAddress("0x5e3Ef299fDDf15eAa0432E6e66473ace8c13D908")
	MainnetEndPoint             = "wss://mainnet.infura.io/ws/v3/86f8d5ba0d524274bce7780a83dbc0a4"
)

func TestSignerToValidatorId(t *testing.T) {
	client, err := NewSimpleClient(MainnetEndPoint)
	if err != nil {
		t.Fatal(err)
	}

	manager, err := StakeManager.NewStakeManager(mainnetStakeManagerContract, client)
	if err != nil {
		t.Fatal(err)
	}

	signer := common.HexToAddress("0x49a7499e26f6311145934163d3aa286aea93fbe9")
	validatorId, err := manager.SignerToValidator(nil, signer)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("validatorId", validatorId) // 115

	valData, err := manager.Validators(nil, validatorId)
	if err != nil {
		t.Fatal(err)
	}

	shareContract := valData.ContractAddress
	t.Log("shareContract", shareContract)

	delay, err := manager.WITHDRAWALDELAY(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("delay", delay)

	interval, err := manager.CheckPointBlockInterval(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("interval", interval)
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
	t.Log(share.Status)          // 1

	valFlag, err := instance.IsValidator(nil, id)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(valFlag)
}

func TestGoerliSigners(t *testing.T) {
	client, err := NewSimpleClient(goerliEndPoint)
	if err != nil {
		t.Fatal(err)
	}

	instance, err := StakeManager.NewStakeManager(goerliStakeManagerContract, client)
	if err != nil {
		t.Fatal(err)
	}

	addr, err := instance.Signers(nil, big.NewInt(1))
	if err != nil {
		t.Fatal(err) //execution reverted
	}
	t.Log(addr)

	id, err := instance.SignerToValidator(nil, addr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("id", id) // 2
	t.Log("id_bytes", hexutil.Encode(id.Bytes()))

	share, err := instance.Validators(nil, id)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("share", share.ContractAddress) // 0x59a82694a675377E010F18F598f5B9dBB83eD968
	t.Log("status", share.Status)         // 1

	valFlag, err := instance.IsValidator(nil, id)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("valFlag", valFlag)
}

func TestMultisigBalance(t *testing.T) {
	client, err := NewSimpleClient(goerliEndPoint)
	if err != nil {
		t.Fatal(err)
	}

	matic, err := MaticToken.NewMaticToken(goerliMaticToken, client)
	if err != nil {
		t.Fatal(err)
	}

	bal, err := matic.BalanceOf(nil, goerliMultisigProxyContract)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bal)
}

func TestTotalStaked(t *testing.T) {
	client, err := NewSimpleClient(goerliEndPoint)
	if err != nil {
		t.Fatal(err)
	}

	manager, err := StakeManager.NewStakeManager(goerliStakeManagerContract, client)
	if err != nil {
		t.Fatal(err)
	}

	share, err := manager.Validators(nil, big.NewInt(9))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("share", share.ContractAddress)

	shr, err := ValidatorShare.NewValidatorShare(share.ContractAddress, client)
	if err != nil {
		t.Fatal(err)
	}

	pool1 := common.HexToAddress("0x03c73f69282e3a1b2a22948bd5a23ce7414490f2")
	total1, _, err := shr.GetTotalStake(nil, pool1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pool1, "totalStake", total1)

	nonce1, err := shr.UnbondNonces(nil, pool1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("nonce1", nonce1)

	matic, err := MaticToken.NewMaticToken(goerliMaticToken, client)
	if err != nil {
		t.Fatal(err)
	}

	bal, err := matic.BalanceOf(nil, pool1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("pool1Bal", bal)

	start, err := manager.Epoch(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("currentEpoch", start)

	delay, err := manager.WITHDRAWALDELAY(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("delay", delay)

	unbond, err := shr.UnbondsNew(nil, pool1, nonce1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("unbond", unbond)

	//interval, err := manager.CheckPointBlockInterval(nil)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log("interval", interval)
	//
	//reward, err := shr.GetLiquidRewards(nil, pool1)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log("reward", reward)
}
