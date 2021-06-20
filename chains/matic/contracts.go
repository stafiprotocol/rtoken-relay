package matic

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stafiprotocol/rtoken-relay/bindings/MaticToken"
	"github.com/stafiprotocol/rtoken-relay/bindings/Multisig"
	"github.com/stafiprotocol/rtoken-relay/bindings/StakeManager"
)

func initStakeManager(stakeManagerCfg interface{}, conn *ethclient.Client) (*StakeManager.StakeManager, error) {
	stakeManagerAddr, ok := stakeManagerCfg.(string)
	if !ok {
		return nil, errors.New("StakeManagerContract not ok")
	}
	return StakeManager.NewStakeManager(common.HexToAddress(stakeManagerAddr), conn)
}

func initMultisig(multisigCfg interface{}, conn *ethclient.Client) (*Multisig.Multisig, error) {
	multisigAddr, ok := multisigCfg.(string)
	if !ok {
		return nil, errors.New("MultisigContract not ok")
	}
	return Multisig.NewMultisig(common.HexToAddress(multisigAddr), conn)
}

func initMaticToken(maticTokenCfg interface{}, conn *ethclient.Client) (*MaticToken.MaticToken, common.Address, error) {
	maticTokenAddr, ok := maticTokenCfg.(string)
	if !ok {
		return nil, ZeroAddress, errors.New("MaticTokenContract not ok")
	}
	addr := common.HexToAddress(maticTokenAddr)
	matic, err := MaticToken.NewMaticToken(addr, conn)
	if err != nil {
		return nil, ZeroAddress, err
	}

	return matic, addr, nil
}

func initMultisend(multiSendCfg interface{}, conn *ethclient.Client) (common.Address, error) {
	multiSendAddr, ok := multiSendCfg.(string)
	if !ok {
		return ZeroAddress, errors.New("MultiSendContract not ok")
	}
	return common.HexToAddress(multiSendAddr), nil
}
