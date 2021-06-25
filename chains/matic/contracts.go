package matic

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stafiprotocol/rtoken-relay/bindings/MaticToken"
	"github.com/stafiprotocol/rtoken-relay/bindings/StakeManager"
)

func initStakeManager(stakeManagerCfg interface{}, conn *ethclient.Client) (*StakeManager.StakeManager, common.Address, error) {
	stakeManagerAddr, ok := stakeManagerCfg.(string)
	if !ok {
		return nil, ZeroAddress, errors.New("StakeManagerContract not ok")
	}
	addr := common.HexToAddress(stakeManagerAddr)
	manager, err := StakeManager.NewStakeManager(common.HexToAddress(stakeManagerAddr), conn)
	if err != nil {
		return nil, ZeroAddress, err
	}

	return manager, addr, nil
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
